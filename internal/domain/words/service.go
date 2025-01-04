package words

import (
	"errors"
	"fmt"

	usersDomain "mono_pardo/internal/domain/users"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"

	"github.com/go-playground/validator"
)

type serviceImpl struct {
	AuthenticationService usersDomain.Service
	Validate              *validator.Validate
	Repository            Repository
}

func NewServiceImpl(
	authenticationService usersDomain.Service,
	validate *validator.Validate,
	repository Repository) Service {
	return &serviceImpl{
		AuthenticationService: authenticationService,
		Validate:              validate,
		Repository:            repository,
	}
}

func (s *serviceImpl) CreateWord(createWordRequest request.CreateWordRequest) error {
	userId, err := s.AuthenticationService.GetUserId(createWordRequest.Token)
	if err != nil {
		return err
	}

	newWord := Word{
		Word:       createWordRequest.Word,
		Definition: createWordRequest.Definition,
		UserId:     userId,
	}

	err = s.Repository.Save(newWord)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) DeleteWord(deleteWordRequest request.DeleteWordRequest) error {
	userId, err := s.AuthenticationService.GetUserId(deleteWordRequest.Token)
	if err != nil {
		return err
	}

	isOwner, err := s.Repository.IsOwnerOfWord(userId, deleteWordRequest.WordId)
	if err != nil {
		return errors.New("cannot check who is owner of the word")
	}

	if isOwner {
		s.Repository.Delete(deleteWordRequest.WordId)
	} else {
		return errors.New("you are not allowed to delete the word")
	}

	return nil
}

func (s *serviceImpl) FindWord(findWordRequest request.FindWordRequest) (response.VocabResponse, error) {
	word, err := s.Repository.FindById(findWordRequest.WordId)
	if err != nil {
		return response.VocabResponse{}, err
	}

	return response.VocabResponse{
		Id:              word.Id,
		Word:            word.Word,
		Definition:      word.Definition,
		CreatedAt:       word.CreatedAt,
		IsLearned:       word.IsLearned,
		Cards:           word.Cards,
		WordTranslation: word.WordTranslation,
		Constructor:     word.Constructor,
		WordAudio:       word.WordAudio,
	}, nil
}

func (s *serviceImpl) GetWords(vocabRequest request.VocabRequest) ([]response.VocabResponse, error) {
	userId, err := s.AuthenticationService.GetUserId(vocabRequest.Token)
	if err != nil {
		return nil, err
	}

	var vocabResponse []response.VocabResponse

	words, words_err := s.Repository.FindByUserId(userId)
	if words_err != nil {
		return nil, words_err
	}

	for _, word := range words {
		vocabResponse = append(vocabResponse, response.VocabResponse{
			Id:              word.Id,
			Word:            word.Word,
			Definition:      word.Definition,
			CreatedAt:       word.CreatedAt,
			IsLearned:       word.IsLearned,
			Cards:           word.Cards,
			WordTranslation: word.WordTranslation,
			Constructor:     word.Constructor,
			WordAudio:       word.WordAudio,
		})
	}

	return vocabResponse, nil
}

func (s *serviceImpl) UpdateWord(updateWordRequest request.UpdateWordRequest) error {
	userId, err := s.AuthenticationService.GetUserId(updateWordRequest.Token)
	if err != nil {
		return err
	}

	trainingFields := map[string]bool{"cards": true, "constructor": true, "word_translation": true, "word_audio": true}

	for _, word := range updateWordRequest.Words {
		if isOwner, err := s.Repository.IsOwnerOfWord(userId, word.WordId); err != nil {
			return fmt.Errorf("error checking ownership of word ID %d: %w", word.WordId, err)
		} else if !isOwner {
			return errors.New("you are not allowed to update the word")
		}

		if err = s.Repository.Update(word); err != nil {
			return err
		}

		for _, update := range word.Updates {
			if trainingFields[update.Field] {
				if err = s.updateWordStatus(word.WordId); err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

func (s *serviceImpl) updateWordStatus(wordId int) error {
	word, err := s.Repository.FindById(wordId)
	if err != nil {
		return fmt.Errorf("cannot find word with id: %d", wordId)
	}

	isLearned := word.Cards && word.WordTranslation && word.Constructor && word.WordAudio

	if word.IsLearned != isLearned {
		req := request.WordUpdate{
			WordId: wordId,
			Updates: []request.FieldUpdate{
				{Field: "is_learned", Value: isLearned},
			},
		}

		if err = s.Repository.Update(req); err != nil {
			return fmt.Errorf("failed to update 'is_learned' status for word ID %d: %w", wordId, err)
		}
	}

	return nil
}

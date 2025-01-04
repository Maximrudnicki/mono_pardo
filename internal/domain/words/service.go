package words

import (
	"errors"

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

	for _, word := range updateWordRequest.Words {
		isOwner, err := s.Repository.IsOwnerOfWord(userId, word.WordId)
		if err != nil {
			return errors.New("cannot check who is owner of the word")
		}

		if isOwner {
			err = s.Repository.Update(word)
			if err != nil {
				return err
			}
		} else {
			return errors.New("you are not allowed to update the word")
		}
	}

	return nil
}

func (s *serviceImpl) UpdateWordStatus(updateWordStatusRequest request.UpdateWordStatusRequest) error {
	userId, err := s.AuthenticationService.GetUserId(updateWordStatusRequest.Token)
	if err != nil {
		return err
	}

	updatedWord := Word{
		Id:        updateWordStatusRequest.WordId,
		IsLearned: updateWordStatusRequest.IsLearned,
		UserId:    userId,
	}

	isOwner, err := s.Repository.IsOwnerOfWord(userId, updateWordStatusRequest.WordId)
	if err != nil {
		return errors.New("cannot check who is owner of the word")
	}

	if isOwner {
		err = s.Repository.UpdateStatus(updatedWord)
		if err != nil {
			return err
		}
	} else {
		return errors.New("you are not allowed to update the word")
	}

	return nil
}

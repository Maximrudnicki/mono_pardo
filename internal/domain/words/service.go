package words

import (
	"errors"

	usersDomain "mono_pardo/internal/domain/users"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"

	"github.com/go-playground/validator"
)

type ServiceImpl struct {
	AuthenticationService usersDomain.Service
	Validate              *validator.Validate
	Repository            Repository
}

func NewServiceImpl(
	authenticationService usersDomain.Service,
	validate *validator.Validate,
	repository Repository) Service {
	return &ServiceImpl{
		AuthenticationService: authenticationService,
		Validate:              validate,
		Repository:            repository,
	}
}

func (s *ServiceImpl) CreateWord(createWordRequest request.CreateWordRequest) error {
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

func (s *ServiceImpl) DeleteWord(deleteWordRequest request.DeleteWordRequest) error {
	userId, err := s.AuthenticationService.GetUserId(deleteWordRequest.Token)
	if err != nil {
		return err
	}

	if isOwner := s.Repository.IsOwnerOfWord(userId, deleteWordRequest.WordId); isOwner {
		s.Repository.Delete(deleteWordRequest.WordId)
	} else {
		return errors.New("you are not allowed to delete the word")
	}

	return nil
}

func (s *ServiceImpl) FindWord(findWordRequest request.FindWordRequest) (response.VocabResponse, error) {
	word, err := s.Repository.FindById(findWordRequest.WordId)
	if err != nil {
		return response.VocabResponse{}, err
	}

	return response.VocabResponse{
		ID:              word.Id,
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

func (s *ServiceImpl) GetWords(vocabRequest request.VocabRequest) ([]response.VocabResponse, error) {
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
			ID:              word.Id,
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

func (s *ServiceImpl) ManageTrainings(manageTrainingsRequest request.ManageTrainingsRequest) error {
	userId, err := s.AuthenticationService.GetUserId(manageTrainingsRequest.Token)
	if err != nil {
		return err
	}

	if isOwner := s.Repository.IsOwnerOfWord(userId, manageTrainingsRequest.WordId); !isOwner {
		return errors.New("you are not allowed to manage trainings for this word")
	}

	err_mt := s.Repository.ManageTrainings(
		manageTrainingsRequest.TrainingResult,
		manageTrainingsRequest.Training,
		manageTrainingsRequest.WordId,
	)

	if err_mt != nil {
		return err_mt
	}

	return nil
}

func (s *ServiceImpl) UpdateWord(updateWordRequest request.UpdateWordRequest) error {
	userId, err := s.AuthenticationService.GetUserId(updateWordRequest.Token)
	if err != nil {
		return err
	}

	updatedWord := Word{
		Id:         updateWordRequest.WordId,
		Definition: updateWordRequest.Definition,
		UserId:     userId,
	}

	if isOwner := s.Repository.IsOwnerOfWord(userId, updateWordRequest.WordId); isOwner {
		err = s.Repository.Update(updatedWord)
		if err != nil {
			return err
		}
	} else {
		return errors.New("you are not allowed to update the word")
	}

	return nil
}

func (s *ServiceImpl) UpdateWordStatus(updateWordStatusRequest request.UpdateWordStatusRequest) error {
	userId, err := s.AuthenticationService.GetUserId(updateWordStatusRequest.Token)
	if err != nil {
		return err
	}

	updatedWord := Word{
		Id:        updateWordStatusRequest.WordId,
		IsLearned: updateWordStatusRequest.IsLearned,
		UserId:    userId,
	}

	if isOwner := s.Repository.IsOwnerOfWord(userId, updateWordStatusRequest.WordId); isOwner {
		err = s.Repository.UpdateStatus(updatedWord)
		if err != nil {
			return err
		}
	} else {
		return errors.New("you are not allowed to update the word")
	}

	return nil
}

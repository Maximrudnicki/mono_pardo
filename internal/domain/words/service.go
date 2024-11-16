package words

import (
	"errors"

	usersDomain "mono_pardo/internal/domain/users"
	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"

	"github.com/go-playground/validator"
)

type VocabService interface {
	CreateWord(createWordRequest request.CreateWordRequest) error
	DeleteWord(deleteWordRequest request.DeleteWordRequest) error
	GetWords(vocabRequest request.VocabRequest) ([]response.VocabResponse, error)
	FindWord(findWordRequest request.FindWordRequest) (response.VocabResponse, error)
	UpdateWord(updateWordRequest request.UpdateWordRequest) error
	UpdateWordStatus(updateWordStatusRequest request.UpdateWordStatusRequest) error
	ManageTrainings(manageTrainingsRequest request.ManageTrainingsRequest) error
	AddWordToStudent(addWordToStudentRequest request.AddWordToStudentRequest) (int, error)
}

type VocabServiceImpl struct {
	AuthenticationService usersDomain.AuthenticationService
	Validate              *validator.Validate
	WordRepository        WordRepository
}

func NewVocabServiceImpl(
	authenticationService usersDomain.AuthenticationService,
	validate *validator.Validate,
	wordRepository WordRepository) VocabService {
	return &VocabServiceImpl{
		AuthenticationService: authenticationService,
		Validate:              validate,
		WordRepository:        wordRepository,
	}
}

// AddWordToStudent implements VocabService.
func (v *VocabServiceImpl) AddWordToStudent(addWordToStudentRequest request.AddWordToStudentRequest) (int, error) {
	newWord := Word{
		Word:       addWordToStudentRequest.Word,
		Definition: addWordToStudentRequest.Definition,
		UserId:     addWordToStudentRequest.UserId,
	}

	wordId, err := v.WordRepository.Add(newWord)
	if err != nil {
		return 0, err
	}

	return wordId, nil
}

// CreateWord implements VocabService.
func (v *VocabServiceImpl) CreateWord(createWordRequest request.CreateWordRequest) error {
	userId, err := v.AuthenticationService.GetUserId(createWordRequest.Token)
	if err != nil {
		return err
	}

	newWord := Word{
		Word:       createWordRequest.Word,
		Definition: createWordRequest.Definition,
		UserId:     userId,
	}

	err = v.WordRepository.Save(newWord)
	if err != nil {
		return err
	}

	return nil
}

// DeleteWord implements VocabService.
func (v *VocabServiceImpl) DeleteWord(deleteWordRequest request.DeleteWordRequest) error {
	userId, err := v.AuthenticationService.GetUserId(deleteWordRequest.Token)
	if err != nil {
		return err
	}

	if isOwner := v.WordRepository.IsOwnerOfWord(userId, deleteWordRequest.WordId); isOwner {
		v.WordRepository.Delete(deleteWordRequest.WordId)
	} else {
		return errors.New("you are not allowed to delete the word")
	}

	return nil
}

// FindWord implements VocabService.
func (v *VocabServiceImpl) FindWord(findWordRequest request.FindWordRequest) (response.VocabResponse, error) {
	word, err := v.WordRepository.FindById(findWordRequest.WordId)
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

// GetWords implements VocabService.
func (v *VocabServiceImpl) GetWords(vocabRequest request.VocabRequest) ([]response.VocabResponse, error) {
	userId, err := v.AuthenticationService.GetUserId(vocabRequest.Token)
	if err != nil {
		return nil, err
	}

	var vocabResponse []response.VocabResponse

	words, words_err := v.WordRepository.FindByUserId(userId)
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

// ManageTrainings implements VocabService.
func (v *VocabServiceImpl) ManageTrainings(manageTrainingsRequest request.ManageTrainingsRequest) error {
	userId, err := v.AuthenticationService.GetUserId(manageTrainingsRequest.Token)
	if err != nil {
		return err
	}

	if isOwner := v.WordRepository.IsOwnerOfWord(userId, manageTrainingsRequest.WordId); !isOwner {
		return errors.New("you are not allowed to manage trainings for this word")
	}

	err_mt := v.WordRepository.ManageTrainings(
		manageTrainingsRequest.TrainingResult,
		manageTrainingsRequest.Training,
		manageTrainingsRequest.WordId,
	)

	if err_mt != nil {
		return err_mt
	}

	return nil
}

// UpdateWord implements VocabService.
func (v *VocabServiceImpl) UpdateWord(updateWordRequest request.UpdateWordRequest) error {
	userId, err := v.AuthenticationService.GetUserId(updateWordRequest.Token)
	if err != nil {
		return err
	}

	updatedWord := Word{
		Id:         updateWordRequest.WordId,
		Definition: updateWordRequest.Definition,
		UserId:     userId,
	}

	if isOwner := v.WordRepository.IsOwnerOfWord(userId, updateWordRequest.WordId); isOwner {
		err = v.WordRepository.Update(updatedWord)
		if err != nil {
			return err
		}
	} else {
		return errors.New("you are not allowed to update the word")
	}

	return nil
}

// UpdateWordStatus implements VocabService.
func (v *VocabServiceImpl) UpdateWordStatus(updateWordStatusRequest request.UpdateWordStatusRequest) error {
	userId, err := v.AuthenticationService.GetUserId(updateWordStatusRequest.Token)
	if err != nil {
		return err
	}

	updatedWord := Word{
		Id:        updateWordStatusRequest.WordId,
		IsLearned: updateWordStatusRequest.IsLearned,
		UserId:    userId,
	}

	if isOwner := v.WordRepository.IsOwnerOfWord(userId, updateWordStatusRequest.WordId); isOwner {
		err = v.WordRepository.UpdateStatus(updatedWord)
		if err != nil {
			return err
		}
	} else {
		return errors.New("you are not allowed to update the word")
	}

	return nil
}

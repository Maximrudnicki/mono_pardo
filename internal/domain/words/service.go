package words

import (
	"fmt"
	"strings"

	"mono_pardo/pkg/data/request"
	"mono_pardo/pkg/data/response"

	"github.com/go-playground/validator"
)

type serviceImpl struct {
	Validate   *validator.Validate
	Repository Repository
}

func NewServiceImpl(validate *validator.Validate, repository Repository) Service {
	return &serviceImpl{
		Validate:   validate,
		Repository: repository,
	}
}

func (s *serviceImpl) CreateWord(createWordRequest request.CreateWordRequest) error {
	newWord, err := NewWord(
		createWordRequest.Word,
		createWordRequest.Definition,
		createWordRequest.UserId,
	)
	if err != nil {
		return fmt.Errorf("invalid word data: %w", err)
	}

	if err = s.Repository.Save(*newWord); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) DeleteWord(deleteWordRequest request.DeleteWordRequest) error {
	if isOwner, err := s.Repository.IsOwnerOfWord(deleteWordRequest.UserId, deleteWordRequest.WordId); err != nil {
		return err
	} else if !isOwner {
		return fmt.Errorf("you are not allowed to delete the word %d", deleteWordRequest.WordId)
	}

	if err := s.Repository.Delete(deleteWordRequest.WordId); err != nil {
		return err
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
	var vocabResponse []response.VocabResponse

	words, err := s.Repository.FindByUserId(vocabRequest.UserId)
	if err != nil {
		return nil, err
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
	trainingFields := map[string]bool{"cards": true, "constructor": true, "word_translation": true, "word_audio": true}

	if err := s.validateWordUpdates(updateWordRequest.Words); err != nil {
		return err
	}

	for _, word := range updateWordRequest.Words {
		if isOwner, err := s.Repository.IsOwnerOfWord(updateWordRequest.UserId, word.WordId); err != nil {
			return err
		} else if !isOwner {
			return fmt.Errorf("you are not allowed to update the word: %d", word.WordId)
		}

		if err := s.Repository.Update(word); err != nil {
			return err
		}

		for _, update := range word.Updates {
			if trainingFields[update.Field] {
				if err := s.updateWordStatus(word.WordId); err != nil {
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
		return err
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
			return fmt.Errorf("failed to update 'is_learned' status for word ID %d", wordId)
		}
	}

	return nil
}

func (s *serviceImpl) validateWordUpdates(updates []request.WordUpdate) error {
	allowedFields := map[string]string{
		"word":             "string",
		"definition":       "string",
		"cards":            "bool",
		"word_translation": "bool",
		"constructor":      "bool",
		"word_audio":       "bool",
	}

	for _, word := range updates {
		if len(word.Updates) == 0 {
			return fmt.Errorf("no updates provided for word ID: %d", word.WordId)
		}

		for _, update := range word.Updates {
			field := strings.TrimSpace(update.Field)
			if field == "" {
				return fmt.Errorf("empty field name not allowed")
			}

			expectedType, validField := allowedFields[field]
			if !validField {
				return fmt.Errorf("invalid field name: %s", field)
			}

			// Type validation based on field
			switch expectedType {
			case "string":
				strValue, ok := update.Value.(string)
				if !ok {
					return fmt.Errorf("field %s requires string value", field)
				}
				if strings.TrimSpace(strValue) == "" {
					return fmt.Errorf("empty value not allowed for field: %s", field)
				}
			case "bool":
				_, ok := update.Value.(bool)
				if !ok {
					return fmt.Errorf("field %s requires boolean value", field)
				}
			}
		}
	}

	return nil
}

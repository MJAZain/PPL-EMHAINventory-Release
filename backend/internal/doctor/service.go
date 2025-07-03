package doctor

import (
	"errors"
	"go-gin-auth/utils"
	"strings"
)

var (
	ErrNotFound     = errors.New("dokter tidak ditemukan")
	ErrInvalidInput = errors.New("input tidak valid atau tidak lengkap")
	ErrSTRExists    = errors.New("nomor STR sudah digunakan oleh dokter lain yang aktif")
)

type Service interface {
	CreateDoctor(doctor *Doctor) (*Doctor, error)
	GetAllDoctors(searchQuery string) ([]Doctor, error)
	GetDoctorByID(id uint) (*Doctor, error)
	UpdateDoctor(id uint, doctor *Doctor) (*Doctor, error)
	DeleteDoctor(id uint) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) encryptDoctorData(d *Doctor) error {
	var err error
	d.FullName, err = utils.Encrypt(d.FullName)
	if err != nil {
		return err
	}
	if d.STRNumber != "" {
		d.STRNumber, err = utils.Encrypt(d.STRNumber)
		if err != nil {
			return err
		}
	}
	d.PhoneNumber, err = utils.Encrypt(d.PhoneNumber)
	if err != nil {
		return err
	}
	if d.Email != "" {
		d.Email, err = utils.Encrypt(d.Email)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) decryptDoctorData(d *Doctor) error {
	var err error
	d.FullName, err = utils.Decrypt(d.FullName)
	if err != nil {
		return err
	}
	if d.STRNumber != "" {
		d.STRNumber, err = utils.Decrypt(d.STRNumber)
		if err != nil {
			return err
		}
	}
	d.PhoneNumber, err = utils.Decrypt(d.PhoneNumber)
	if err != nil {
		return err
	}
	if d.Email != "" {
		d.Email, err = utils.Decrypt(d.Email)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) CreateDoctor(doctor *Doctor) (*Doctor, error) {
	doctor.FullName = strings.TrimSpace(doctor.FullName)
	if doctor.FullName == "" || doctor.Specialization == "" || doctor.PhoneNumber == "" {
		return nil, ErrInvalidInput
	}

	if doctor.STRNumber != "" {
		allActiveDoctors, err := s.repository.GetAllActive()
		if err != nil {
			return nil, err
		}
		for _, activeDoc := range allActiveDoctors {
			decryptedSTR, _ := utils.Decrypt(activeDoc.STRNumber)
			if decryptedSTR == doctor.STRNumber {
				return nil, ErrSTRExists
			}
		}
	}

	if err := s.encryptDoctorData(doctor); err != nil {
		return nil, errors.New("gagal mengenkripsi data dokter")
	}

	if doctor.Status == "" {
		doctor.Status = "Aktif"
	}

	newDoctor, err := s.repository.Create(doctor)
	if err != nil {
		return nil, err
	}

	if err := s.decryptDoctorData(newDoctor); err != nil {
		return nil, errors.New("gagal mendekripsi data untuk respons")
	}
	return newDoctor, nil
}

func (s *service) GetAllDoctors(searchQuery string) ([]Doctor, error) {
	doctors, err := s.repository.GetAll(searchQuery)
	if err != nil {
		return nil, err
	}

	decryptedDoctors := make([]Doctor, 0)
	for _, doc := range doctors {
		if err := s.decryptDoctorData(&doc); err == nil {
			decryptedDoctors = append(decryptedDoctors, doc)
		}
	}

	if searchQuery != "" {
		var filtered []Doctor
		searchQuery = strings.ToLower(searchQuery)
		for _, d := range decryptedDoctors {
			if strings.Contains(strings.ToLower(d.FullName), searchQuery) || strings.Contains(d.STRNumber, searchQuery) {
				filtered = append(filtered, d)
			}
		}
		return filtered, nil
	}
	return decryptedDoctors, nil
}

func (s *service) GetDoctorByID(id uint) (*Doctor, error) {
	doctor, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.decryptDoctorData(doctor); err != nil {
		return nil, errors.New("gagal mendekripsi data dokter")
	}
	return doctor, nil
}

func (s *service) UpdateDoctor(id uint, input *Doctor) (*Doctor, error) {
	_, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.encryptDoctorData(input); err != nil {
		return nil, errors.New("gagal mengenkripsi data untuk update")
	}

	updatedDoctor, err := s.repository.Update(id, input)
	if err != nil {
		return nil, err
	}

	if err := s.decryptDoctorData(updatedDoctor); err != nil {
		return nil, errors.New("gagal mendekripsi data untuk respons")
	}
	return updatedDoctor, nil
}

func (s *service) DeleteDoctor(id uint) error {
	if _, err := s.repository.GetByID(id); err != nil {
		return err
	}
	return s.repository.Delete(id)
}

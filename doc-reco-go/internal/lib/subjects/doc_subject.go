package subjects

import (
	"doc-reco-go/internal/utils"
	"fmt"
)

type DocSubject struct {
	SubjectId      int
	Slug           string
	Name           string
	IconIdentifier string
	Type           string

	subjects        []DocSubject
	masterSubject   map[int]string
	masterSubjectId []int
	docMasterMap    map[int][]int
}

func (s DocSubject) Init(subjects []DocSubject, masterSubject map[int]string, docMasterMap map[int][]int) DocSubject {
	if subjects == nil {
		s.subjects = SUBJECT_LIST
	} else {
		s.subjects = subjects
	}
	if masterSubject == nil {
		s.masterSubject = MASTER_SUBJECTS_MAP
	} else {
		s.masterSubject = masterSubject
	}
	if docMasterMap == nil {
		s.docMasterMap = MAP_DOC_TO_MASTER_SUBJECT
	} else {
		s.docMasterMap = docMasterMap
	}
	if s.IconIdentifier == "" {
		s.IconIdentifier = "regular"
	}

	for id, _ := range s.masterSubject {
		s.masterSubjectId = append(s.masterSubjectId, id)
	}
	return s
}

func (s DocSubject) GetSubjectIdFilter(klass , docSubjectId int) ([]int, error) {
	if docSubjectId == 0 {
		return nil, nil
	}
	return s.getMasterSubjectsByKlass(klass, docSubjectId)
}

func (s DocSubject) getMasterSubjectsByKlass(klass , docSubjectId int) ([]int, error) {
	if docSubjectId == -1{
		if klass < 0 || klass > 13 {
			return nil, fmt.Errorf("invalid class value: %d", klass)
		} else if klass <= 10 {
			return s.getOtherMasterSubjectIds(LT_CLASS_10_SUBJECTS), nil
		} else if klass < 13 {
			return s.getOtherMasterSubjectIds(GT_CLASS_10_SUBJECTS), nil
		} else {
			return s.getOtherMasterSubjectIds(NO_CLASS_SUBJECTS), nil
		}
	}
	val, ok := s.docMasterMap[docSubjectId]
	if !ok	{
		return nil, fmt.Errorf("master subject mapping not found for id: %d", docSubjectId)
	}
	return val, nil
}

func (s DocSubject) getSubjectRespFromSlugList(slugList []string) ([]DocSubject, []DocSubject) {
	var subjectList, otherSubjectList []DocSubject
	for _, sub := range s.subjects {
		if utils.FindElementInSlice(sub.Slug, slugList) {
			subjectList = append(subjectList, sub)
		} else if sub.SubjectId != -1 {
			otherSubjectList = append(otherSubjectList, sub)
		}
	}
	return subjectList, otherSubjectList
}

func (s DocSubject) getOtherMasterSubjectIds(slugList []string) []int {
	subjectList, _ := s.getSubjectRespFromSlugList(slugList)

	var masterSubjectIds []int
	for _, sub := range subjectList {
		if sub.Slug != "other" {
			masterSubjectIds = append(masterSubjectIds, s.docMasterMap[sub.SubjectId]...)
		}
	}

	otherMasterSubjectIds := utils.GetSliceDifference(s.masterSubjectId, masterSubjectIds)
	return otherMasterSubjectIds
}
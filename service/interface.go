package service

import (
	"sync"

	"github.com/ThreeCatsLoveFish/medalhelper/dto"
)

type IConcurrency interface {
	// Exec the action of child and execute retry backup if
	Exec(user User, work *sync.WaitGroup, child IExec) []dto.MedalInfo
}

type IExec interface {
	// Do represent real action
	Do(user User, medal dto.MedalInfo) bool
	// Finish represent action complete
	Finish(user User, medal []dto.MedalInfo)
}

// Action represent a single action for a single user
type IAction interface {
	IConcurrency
	IExec
}

package anthive

import (
	"errors"
	"strings"
)

// Modes for FieldInfo
const (
	FIELD_ANTS  = iota // On Reading Ants
	FIELD_ROOMS        // On Reading Rooms
	FIELD_PATHS        // On Reading Paths | Relations
)

// Modes for Path
const (
	REVERSED = -1 // directed, REVERSED path (from end to start)
	BLOCKED  = 0  // blocked path (from start to end)
	STABLE   = 1  // double directed path
)

// The found paths are saved in Result. Using for write result to writer
type Result struct {
	AntsCount int
	Paths     []*list
}

// Stores information about the graph, the data being read, and the result. Using for find paths
type anthive struct {
	// For Reading Data
	FieldInfo *fieldInfo
	// Main Data
	AntsCount  int
	Start, End string
	Rooms      map[string]*room
	// Results
	StepsCount int
	Result     *Result
}

type room struct {
	Name                string        // Name
	X, Y                int           // Coordinates
	Paths               map[*room]int // Paths with state -> REVERSED || BLOCKED || STABLE
	ParentIn, ParentOut *room         // Store parents for new path
	VisitIn, VisitOut   bool          // Flag for checking while traversing
	Weight              [2]int        // Out weight in 0 index, In in 1
	Separated           bool          // Flag for checking separated node
}

// with fieldInfo, we understand What data we fill in for the anthive
type fieldInfo struct {
	MODE             byte                 // FIELD_ANTS | FIELD_ROOMS | FIELD_PATHS
	Start, End       bool                 // Should Be True
	IsStart, IsEnd   bool                 // For Know Which Room is Reading
	UsingCoordinates map[int]map[int]bool // Chekking for unique Coordinates on Rooms
}

// List of node which has room.
type node struct {
	Room *room
	Next *node
}

// List of Room nodes. Used to store found paths
type list struct {
	Len   int
	Front *node
	Back  *node
}

type antStruct struct {
	Num  int
	Path int
	Pos  int
	Next *antStruct
}

// queue of ants. Used for write ant position on every step
type antQueue struct {
	Front *antStruct
	Back  *antStruct
}

// for sorting rooms in queue
type weightNode struct {
	Room   *room
	Weight int
	Mark   bool // false if it's in_node, true if it's out_node
	Next   *weightNode
}

// queue for sorted rooms by weight
type sortedQueue struct {
	Front *weightNode
	Back  *weightNode
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Createanthive - returns anthive by default data
func Createanthive() *anthive {
	result := &anthive{}
	result.Rooms = make(map[string]*room)
	result.FieldInfo = &fieldInfo{UsingCoordinates: make(map[int]map[int]bool)}
	result.Result = &Result{}
	return result
}

// ValidateByFieldInfo - returns an error if something was missed by the scanner
func (a *anthive) ValidateByFieldInfo() error {
	if a.FieldInfo.MODE != FIELD_PATHS {
		switch a.FieldInfo.MODE {
		case FIELD_ANTS:
			return errors.New("here is no Ants")
		case FIELD_ROOMS:
			return errors.New("here is no Rooms or Paths")
		default:
			return errors.New("func Validate returns error")
		}
	} else {
		if !a.FieldInfo.Start {
			return errors.New("please set ##start room")
		} else if !a.FieldInfo.End {
			return errors.New("please set ##end room")
		}
	}
	return nil
}

// ReadDataFromLine - reading the line, it replenishes the data about the anthive. (FieldInfo understands what the string is)
func (a *anthive) ReadDataFromLine(line string) error {
	if line == "" || strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
		return nil
	}
	switch a.FieldInfo.MODE {
	case FIELD_PATHS:
		err := a.SetPathsFromLine(line)
		if err != nil {
			return err
		}
	case FIELD_ROOMS:
		if strings.HasPrefix(line, "##") {
			if line == "##start" && !a.FieldInfo.Start && !a.FieldInfo.IsStart && !a.FieldInfo.IsEnd {
				a.FieldInfo.IsStart = true
				return nil
			} else if line == "##end" && !a.FieldInfo.End && !a.FieldInfo.IsEnd && !a.FieldInfo.IsStart {
				a.FieldInfo.IsEnd = true
				return nil
			}
			return errors.New("error with ## command")
		}
		if a.FieldInfo.IsStart || a.FieldInfo.IsEnd {
			err := a.SetMainRooms(line, a.FieldInfo.IsStart)
			if err != nil {
				return err
			}
			if a.FieldInfo.IsStart {
				a.FieldInfo.IsStart = false
				a.FieldInfo.Start = true
			} else {
				a.FieldInfo.IsEnd = false
				a.FieldInfo.End = true
			}
			return err
		} else if len(strings.Split(line, " ")) != 3 {
			a.FieldInfo.MODE = FIELD_PATHS
			a.FieldInfo.UsingCoordinates = nil
			return a.ReadDataFromLine(line)
		} else {
			_, err := a.SetRoomFromLine(line)
			return err
		}
	case FIELD_ANTS:
		err := a.SetAntsFromLine(line)
		if err != nil {
			return err
		}
		a.FieldInfo.MODE = FIELD_ROOMS
	}
	return nil
}

// Match - Finds paths, returns an error if it does not find a single path. Paths are saved in anthive.Result
func (a *anthive) Match() error {
	for {
		if !searchShortPath(a) {
			// path not found, then check for prev path count
			if a.StepsCount > 0 {
				return nil
			}
			return errors.New("path not found")
		}
		if !checkEffective(a) {
			return nil
		}
	}
}

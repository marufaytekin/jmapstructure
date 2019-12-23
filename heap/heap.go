package heap

import (
	"github.com/mitchellh/mapstructure"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Stats struct {
	Regions  int64
	Capacity int64
	Used     int64
	Free     int64
}

type Config struct {
	MinHeapFreeRatio         int64
	MaxHeapFreeRatio         int64
	MaxHeapSize              int64
	NewSize                  int64
	MaxNewSize               int64
	OldSize                  int64
	NewRatio                 int64
	SurvivorRatio            int64
	MetaspaceSize            int64
	CompressedClassSpaceSize int64
	MaxMetaspaceSize         int64
	G1HeapRegionSize         int64
}

type Heap struct {
	JavaVersion     string
	HeapConfig      Config
	G1heap          Stats
	Edenspace       Stats
	Survuvorspace   Stats
	G1oldgeneration Stats
}

type Map struct {
	JavaVersion     string
	HeapConfig      map[string]int64
	G1heap          map[string]int64
	Edenspace       map[string]int64
	Survivorspace   map[string]int64
	G1oldgeneration map[string]int64
}

func newHeapMap() *Map {
	m := &Map{
		JavaVersion:     "",
		HeapConfig:      make(map[string]int64),
		G1heap:          make(map[string]int64),
		Edenspace:       make(map[string]int64),
		Survivorspace:   make(map[string]int64),
		G1oldgeneration: make(map[string]int64),
	}
	return m
}

// this function tries jmap command first to dump the heap stats.
// if it fails it tries jhsdb.
func Get(javahome string, pid string) (*Heap, error) {
	m := newHeapMap()
	c := javahome + "/bin/jmap"
	args := []string{"-heap", pid}
	h, err := m.getCurrentHeap(c, args)
	if err == nil {
		return h, err
	}
	c = javahome + "/bin/jhsdb"
	args = []string{"jmap", "--heap", "--pid", pid}

	return m.getCurrentHeap(c, args)
}

func (m *Map) getCurrentHeap (c string, args []string) (*Heap, error) {
	o, err := m.runJmapCmd(c, args)
	if err != nil {
		return nil, err
	}
	err = m.parseJmapOutput(o)
	if err != nil {
		return nil, err
	}

	return m.toStruct()
}

func (m *Map) runJmapCmd(c string, args []string) (string, error) {

	cmd := exec.Command(c, args...)
	o, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(o), nil
}

func (m *Map) parseJmapOutput (o string) error {
	// starting with heap config
	curr := m.HeapConfig
	for _, line := range strings.Split(string(o), "\n") {
		if strings.Contains(line, "JVM version is") {
			m.JavaVersion = parseLine(line)[3]
		}
		currLine := parseLine(line)
		if strings.Contains(line, "G1 Heap:") {
			curr = m.G1heap
		}
		if strings.Contains(line, "Eden Space:") {
			curr = m.Edenspace
		}
		if strings.Contains(line, "Survivor Space:") {
			curr = m.Survivorspace
		}
		if strings.Contains(line, "G1 Old Generation:") {
			curr = m.G1oldgeneration
		}
		if strings.Contains(line, "=") {
			val, err := strconv.ParseInt(currLine[2], 10, 64)
			if err != nil {
				return err
			}
			curr[currLine[0]] = val
		}
	}
	return nil
}

func (m *Map) toStruct() (*Heap, error) {
	conf := &Config{}
	g1h  := &Stats{}
	es   := &Stats{}
	ss   := &Stats{}
	g1o  := &Stats{}
	err := mapstructure.Decode(m.HeapConfig, conf)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(m.G1heap, g1h)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(m.Edenspace, es)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(m.Survivorspace, ss)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(m.G1oldgeneration, g1o)
	if err != nil {
		return nil, err
	}
	h := &Heap{
		JavaVersion:     m.JavaVersion,
		HeapConfig:      *conf,
		G1heap:          *g1h,
		Edenspace:       *es,
		Survuvorspace:   *ss,
		G1oldgeneration: *g1o,
	}
	return h, nil
}

func parseLine(line string) []string {
	space := regexp.MustCompile(`\s+`)
	s := strings.Trim(space.ReplaceAllString(line, " "), " ")
	return strings.Split(s, " ")
}

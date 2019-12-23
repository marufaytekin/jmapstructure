package heap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	output = `Attaching to process ID 23951, please wait...
Debugger attached successfully.
Server compiler detected.
JVM version is 12.0.2+10

using thread-local object allocation.
Garbage-First (G1) GC with 13 thread(s)

Heap Configuration:
   MinHeapFreeRatio         = 40
   MaxHeapFreeRatio         = 70
   MaxHeapSize              = 536870912 (512.0MB)
   NewSize                  = 1363144 (1.2999954223632812MB)
   MaxNewSize               = 321912832 (307.0MB)
   OldSize                  = 5452592 (5.1999969482421875MB)
   NewRatio                 = 2
   SurvivorRatio            = 8
   MetaspaceSize            = 21807104 (20.796875MB)
   CompressedClassSpaceSize = 1073741824 (1024.0MB)
   MaxMetaspaceSize         = 17592186044415 MB
   G1HeapRegionSize         = 1048576 (1.0MB)

Heap Usage:
G1 Heap:
   regions  = 512
   capacity = 536870912 (512.0MB)
   used     = 18579456 (17.71875MB)
   free     = 518291456 (494.28125MB)
   3.460693359375% used
G1 Young Generation:
Eden Space:
   regions  = 17
   capacity = 28311552 (27.0MB)
   used     = 17825792 (17.0MB)
   free     = 10485760 (10.0MB)
   62.96296296296296% used
Survivor Space:
   regions  = 0
   capacity = 0 (0.0MB)
   used     = 0 (0.0MB)
   free     = 0 (0.0MB)
   0.0% used
G1 Old Generation:
   regions  = 3
   capacity = 508559360 (485.0MB)
   used     = 753664 (0.71875MB)
   free     = 507805696 (484.28125MB)
   0.1481958762886598% used
`
)

func TestParseJmapOutput(t *testing.T) {
	m := newHeapMap()

	err := m.parseJmapOutput(output)
	assert.NoError(t, err)

	h, err := m.toStruct()
	assert.NoError(t, err)

	assert.Equal(t, "12.0.2+10", h.JavaVersion)
	assert.Equal(t, int64(40), h.HeapConfig.MinHeapFreeRatio)
	assert.Equal(t,int64(70), h.HeapConfig.MaxHeapFreeRatio)
	assert.Equal(t, int64(536870912), h.HeapConfig.MaxHeapSize)
	assert.Equal(t,int64(1363144), h.HeapConfig.NewSize)
	assert.Equal(t, int64(321912832), h.HeapConfig.MaxNewSize)
	assert.Equal(t, int64(5452592), h.HeapConfig.OldSize)
	assert.Equal(t, int64(2), h.HeapConfig.NewRatio)
	assert.Equal(t,int64(8), h.HeapConfig.SurvivorRatio)
	assert.Equal(t,int64(21807104), h.HeapConfig.MetaspaceSize)
	assert.Equal(t,int64(1073741824), h.HeapConfig.CompressedClassSpaceSize)
	assert.Equal(t, int64(17592186044415), h.HeapConfig.MaxMetaspaceSize)
	assert.Equal(t,int64(1048576), h.HeapConfig.G1HeapRegionSize)

	assert.Equal(t, int64(512), h.G1heap.Regions)
	assert.Equal(t, int64(536870912), h.G1heap.Capacity)
	assert.Equal(t, int64(18579456), h.G1heap.Used)
	assert.Equal(t, int64(518291456), h.G1heap.Free)

	assert.Equal(t, int64(17), h.Edenspace.Regions)
	assert.Equal(t, int64(28311552), h.Edenspace.Capacity)
	assert.Equal(t, int64(17825792), h.Edenspace.Used)
	assert.Equal(t, int64(10485760), h.Edenspace.Free)

	assert.Equal(t, int64(0), h.Survuvorspace.Regions)
	assert.Equal(t, int64(0), h.Survuvorspace.Capacity)
	assert.Equal(t, int64(0), h.Survuvorspace.Used)
	assert.Equal(t, int64(0), h.Survuvorspace.Free)

	assert.Equal(t, int64(3), h.G1oldgeneration.Regions)
	assert.Equal(t, int64(508559360), h.G1oldgeneration.Capacity)
	assert.Equal(t, int64(753664), h.G1oldgeneration.Used)
	assert.Equal(t, int64(507805696), h.G1oldgeneration.Free)

}
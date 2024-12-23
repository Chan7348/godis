package timewheel

import (
	"container/list"
	"time"

	"github.com/Chan7348/godis/lib/logger"
)

type TimeWheel struct {
	interval time.Duration
	ticker   *time.Ticker
	slots    []*list.List

	keyToLocation map[string]*location

	currentSlot int
	totalSlots  int

	addTaskChannel    chan task
	removeTaskChannel chan string
	stopChannel       chan bool
}

type location struct {
	slot  int
	etask *list.Element
}

type task struct {
	delay  time.Duration
	circle int
	key    string
	job    func()
}

func New(interval time.Duration, totalSlots int) *TimeWheel {
	timeWheel := &TimeWheel{
		interval:          interval,
		slots:             make([]*list.List, totalSlots),
		keyToLocation:     make(map[string]*location),
		currentSlot:       0,
		totalSlots:        totalSlots,
		addTaskChannel:    make(chan task),
		removeTaskChannel: make(chan string),
		stopChannel:       make(chan bool),
	}
	timeWheel.initSlots()
	return timeWheel
}

func (tw *TimeWheel) initSlots() {
	for i := range tw.totalSlots {
		tw.slots[i] = list.New()
	}
}

func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
}

func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

func (tw *TimeWheel) AddJob(delay time.Duration, key string, job func()) {
	if delay < 0 {
		return
	}
	tw.addTaskChannel <- task{delay: delay, key: key, job: job}
}

func (tw *TimeWheel) RemoveJob(key string) {
	if key == "" {
		return
	}
	tw.removeTaskChannel <- key
}

func (tw *TimeWheel) run() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(&task)
		case key := <-tw.removeTaskChannel:
			tw.removeTask(key)
		case <-tw.stopChannel:
			tw.Stop()
			return
		}
	}
}

func (tw *TimeWheel) tickHandler() {
	// get list
	list := tw.slots[tw.currentSlot]
	if tw.currentSlot == tw.totalSlots-1 {
		tw.currentSlot = 0
	} else {
		tw.currentSlot++
	}
	go tw.scanAndRunTask(list)
}

func (tw *TimeWheel) scanAndRunTask(list *list.List) {
	for e := list.Front(); e != nil; {
		task := e.Value.(*task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(err)
				}
			}()
			task.job()
		}()
		next := e.Next()
		list.Remove(e)
		if task.key != "" {
			delete(tw.keyToLocation, task.key)
		}
		e = next
	}
}

func (tw *TimeWheel) addTask(task *task) {
	position, circle := tw.getPositionAndCircle(task.delay)
	task.circle = circle

	e := tw.slots[position].PushBack(task)
	location := &location{
		slot:  position,
		etask: e,
	}
	if task.key != "" {
		_, ok := tw.keyToLocation[task.key]
		if ok {
			tw.removeTask(task.key)
		}
	}
	tw.keyToLocation[task.key] = location
}

func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delayInSeconds := int(d.Seconds())
	intervalInSeconds := int(tw.interval.Seconds())

	circle = int(delayInSeconds / intervalInSeconds / tw.totalSlots)
	pos = int(tw.currentSlot+delayInSeconds/intervalInSeconds) % tw.totalSlots

	return
}

func (tw *TimeWheel) removeTask(key string) {
	location, ok := tw.keyToLocation[key]
	if !ok {
		return
	}
	list := tw.slots[location.slot]
	list.Remove(location.etask)   // 从链表中删除节点
	delete(tw.keyToLocation, key) // 从mapping中删除key的映射
}

package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Реализовать функцию, которая будет объединять один или более done-каналов в
single-канал, если один из его составляющих каналов закроется.

Очевидным вариантом решения могло бы стать выражение при использовании select, которое бы
реализовывало эту связь, однако иногда неизвестно общее число done-каналов, с
которыми вы работаете в рантайме. В этом случае удобнее использовать вызов
единственной функции, которая, приняв на вход один или более or-каналов,
реализовывала бы весь функционал.
*/

var or func(channels ...<-chan interface{}) <-chan interface{}

func init() {
	or = MergeDoneChannels
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v", time.Since(start))
}

func MergeDoneChannels(channels ...<-chan interface{}) <-chan interface{} {
	// результирующий канал
	collectorCh := make(chan interface{})

	// done канал для прекращения прослушивания input каналов горутинами
	broadcastCh := make(chan struct{})

	// wg для того, чтобы гарантировать, что collectorCh закроется только после того как все горутины слушающие
	// input каналы завершаться, тогда не будет ситуации, что канал collectorCh закроется более одного раза
	var wg sync.WaitGroup
	wg.Add(len(channels))

	// На каждый input канал создается своя горутина, которая ее обрабатывает, если done, то сообщается об этом
	// всем остальным горутинам и они завершают работу декрементируя wg.
	for _, ch := range channels {
		go func(c <-chan interface{}) {
			select {
			case <-c:
				// Пришло сообщение о done -> сообщаем остальным горутинам через broadcastCh и прекращаем горутину
				close(broadcastCh)
				wg.Done()
				break
			case <-broadcastCh:
				// Кому-то из остальных каналов пришло done -> завершаем обработку канала и "убиваем" горутину
				wg.Done()
				break
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(collectorCh)
	}()

	return collectorCh
}

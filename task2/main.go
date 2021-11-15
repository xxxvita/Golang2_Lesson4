/*
ps -aux
kill -15 PID
*/
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGTERM)

	ctx, funcCLose := context.WithCancel(context.Background())

	// Обработка сигналов ОС
	go func(ctx context.Context) {
		_ = <-chanSignal
		log.Println("Пришёл сигнал SIGTERM")

		ctxClose, funcClose := context.WithTimeout(ctx, time.Second)
		defer funcClose()

		select {
		case <-ctxClose.Done():
			// Завершается main-context
			funcCLose()
		}
	}(ctx)

	// Запуска долгоработающей функции
	err := startService(ctx)
	if err != nil {
		log.Printf("Ошибка при выполнении сервиса %s", err)
		return
	}

	select {
	case <-ctx.Done():
		log.Println("Все потоки завершены")
	}
}

func startService(ctx context.Context) error {
	ctx1, cancelFunc := context.WithTimeout(ctx, 1*time.Minute*5)
	defer cancelFunc()

	// Горунтина, которая симулирует долгую работу и всё время работает
	// Должна завершиться не поздее чем через секунду после сигнала SIGTERM
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Service aborted")
				return
			default:
				time.Sleep(time.Second)
				log.Printf("Service is work")
			}
		}
	}(ctx1)

	var err error
	select {
	case <-ctx1.Done():
		err = ctx1.Err()
	}

	return err
}

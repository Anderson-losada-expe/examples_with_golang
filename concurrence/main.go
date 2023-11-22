package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Account struct {
	balance int64
}

func (a *Account) getBalance() int64 {
	return atomic.LoadInt64(&a.balance)
}

func (a *Account) deposit(amount int64) {
	atomic.AddInt64(&a.balance, amount)
}

func (a *Account) withdraw(amount int64) bool {
	for {
		currentBalance := atomic.LoadInt64(&a.balance)
		if currentBalance < amount {
			return false
		}

		newBalance := currentBalance - amount
		if atomic.CompareAndSwapInt64(&a.balance, currentBalance, newBalance) {
			return true
		}
	}
}

func transfer(sender, receiver *Account, amount int64, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()

	fmt.Printf("Transferencia iniciada. Saldo antes de la transferencia: %d\n", sender.getBalance())

	// Simular operación compleja
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

	if sender.withdraw(amount) {
		receiver.deposit(amount)
		fmt.Printf("Transferencia completada. Saldo después de la transferencia: %d\n", sender.getBalance())
	} else {
		errCh <- fmt.Errorf("fondos insuficientes para la transferencia")
	}
}

func main() {
	var wg sync.WaitGroup
	errCh := make(chan error, 10)

	// Crear cuentas
	account1 := &Account{balance: 50000}
	account2 := &Account{balance: 50000}

	// Simular múltiples transferencias
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go transfer(account1, account2, int64(rand.Intn(100)), &wg, errCh)
	}

	// Esperar a que todas las transferencias terminen
	wg.Wait()

	// Cerrar el canal de errores
	close(errCh)

	// Manejar errores
	for err := range errCh {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Printf("Saldo final de la cuenta 1: %d\n", account1.getBalance())
	fmt.Printf("Saldo final de la cuenta 2: %d\n", account2.getBalance())
}

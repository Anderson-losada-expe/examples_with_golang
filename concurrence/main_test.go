package main

import (
	"sync"
	"testing"
)

func TestAccountOperations(t *testing.T) {
	// Crear una cuenta
	account := &Account{balance: 50000}

	// Probar getBalance
	if balance := account.getBalance(); balance != 50000 {
		t.Errorf("Error en getBalance. Se esperaba 50000, se obtuvo %d", balance)
	}

	// Probar deposit
	account.deposit(10000)
	if balance := account.getBalance(); balance != 60000 {
		t.Errorf("Error en deposit. Se esperaba 60000, se obtuvo %d", balance)
	}

	// Probar withdraw
	success := account.withdraw(30000)
	if !success {
		t.Error("Error en withdraw. Se esperaba éxito, se obtuvo fallo")
	}
	if balance := account.getBalance(); balance != 30000 {
		t.Errorf("Error en withdraw. Se esperaba 30000, se obtuvo %d", balance)
	}

	// Probar withdraw con fondos insuficientes
	success = account.withdraw(50000)
	if success {
		t.Error("Error en withdraw. Se esperaba fallo, se obtuvo éxito")
	}
	if balance := account.getBalance(); balance != 30000 {
		t.Errorf("Error en withdraw. Se esperaba 30000, se obtuvo %d", balance)
	}
}

func TestTransfer(t *testing.T) {
	// Crear cuentas
	account1 := &Account{balance: 50000}
	account2 := &Account{balance: 50000}

	// Transferir una cantidad válida
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	wg.Add(1)
	go transfer(account1, account2, 10000, &wg, errCh)
	wg.Wait()

	// Verificar saldos finales
	if balance := account1.getBalance(); balance != 40000 {
		t.Errorf("Error en TestTransfer. Saldo final de la cuenta 1 incorrecto. Se esperaba 40000, se obtuvo %d", balance)
	}

	if balance := account2.getBalance(); balance != 60000 {
		t.Errorf("Error en TestTransfer. Saldo final de la cuenta 2 incorrecto. Se esperaba 60000, se obtuvo %d", balance)
	}

	// Transferir una cantidad inválida
	wg.Add(1)
	go transfer(account1, account2, 50000, &wg, errCh)
	wg.Wait()

	// Verificar error
	select {
	case err := <-errCh:
		if err == nil {
			t.Error("Error en TestTransfer. Se esperaba un error, pero no se recibió ninguno")
		}
	default:
		t.Error("Error en TestTransfer. No se recibió ningún error, pero se esperaba uno")
	}
}

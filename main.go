package main

import (
	"fmt"
	"math/big"
)

// Упростим задачу, достаточно посчитать количество вариантов окончания на ноль в последовательности
// Например у последовательности 1 < n < 100 два варианта, 0 нулей в конце (у чисел 2, 3, 4 и тд)
// и 1 ноль в конце (10, 20 и тд). Это равносильно работе функции g (по количеству возможных результатов)

// Подсчитывает количество нулей в конце числа
func countTrailingZeros(n *big.Int) int {
	if n.Cmp(big.NewInt(0)) == 0 {
		return 1
	}

	count := 0
	temp := new(big.Int).Set(n)
	ten := big.NewInt(10)
	remainder := new(big.Int)

	for {
		temp.DivMod(temp, ten, remainder)
		if remainder.Cmp(big.NewInt(0)) != 0 {
			break
		}
		count++
	}

	return count
}

// Для больших чисел в заданном диапазоне
func countZeroEndingCombinations(start, end *big.Int) int {
	if start.Cmp(end) > 0 || end.Cmp(big.NewInt(0)) <= 0 {
		return 0
	}

	combinations := make(map[int]bool)

	// Находим максимальное количество нулей, которое может быть в диапазоне
	endStr := end.String()
	maxPossibleZeros := len(endStr) - 1

	// Проверяем все возможные количества нулей от 0 до максимального
	for zeros := 0; zeros <= maxPossibleZeros; zeros++ {
		if hasNumbersWithTrailingZeros(start, end, zeros) {
			combinations[zeros] = true
		}
	}

	return len(combinations)
}

// Проверяет, есть ли в диапазоне числа с определенным количеством нулей в конце
func hasNumbersWithTrailingZeros(start, end *big.Int, targetZeros int) bool {
	if targetZeros == 0 {
		// Для нуля нулей - ищем любое число, не оканчивающееся на ноль
		return hasNumbersNotEndingWithZero(start, end)
	}

	// Для чисел с нулями - ищем числа вида k * 10^targetZeros
	power := new(big.Int)
	power.Exp(big.NewInt(10), big.NewInt(int64(targetZeros)), nil)

	// Находим первое число в диапазоне, которое кратно 10^targetZeros
	startDivided := new(big.Int).Div(start, power)
	if new(big.Int).Mul(startDivided, power).Cmp(start) < 0 {
		startDivided.Add(startDivided, big.NewInt(1))
	}

	// Проверяем числа вида k * 10^targetZeros
	multiplier := new(big.Int).Set(startDivided)
	candidate := new(big.Int)

	for {
		candidate.Mul(multiplier, power)

		if candidate.Cmp(end) > 0 {
			break
		}

		if candidate.Cmp(start) >= 0 {
			// Проверяем, что у числа именно targetZeros нулей, а не больше
			actualZeros := countTrailingZeros(candidate)
			if actualZeros == targetZeros {
				return true
			}
		}

		multiplier.Add(multiplier, big.NewInt(1))
	}

	return false
}

// Проверяет, есть ли в диапазоне числа, не оканчивающиеся на ноль
func hasNumbersNotEndingWithZero(start, end *big.Int) bool {
	current := new(big.Int).Set(start)

	// Для небольших диапазонов проверяем все числа
	rangeSize := new(big.Int).Sub(end, start)
	if rangeSize.Cmp(big.NewInt(1000)) <= 0 {
		for current.Cmp(end) <= 0 {
			if countTrailingZeros(current) == 0 {
				return true
			}
			current.Add(current, big.NewInt(1))
		}
		return false
	}

	// Для больших диапазонов используем математическую логику
	// В любом диапазоне из 10 последовательных чисел только одно заканчивается на 0
	// Если диапазон содержит больше 10 чисел, то точно есть числа без нулей
	if rangeSize.Cmp(big.NewInt(10)) >= 0 {
		return true
	}

	// Для диапазонов 1-10 чисел проверяем каждое
	for current.Cmp(end) <= 0 {
		if countTrailingZeros(current) == 0 {
			return true
		}
		current.Add(current, big.NewInt(1))
	}

	return false
}

func main() {
	// Задаем начало и конец диапазона
	start := new(big.Int)
	end := new(big.Int)

	// Последовательность от 2 до 999999999999999999999999999999 (1 < n < 10^30)
	start.SetString("2", 10)
	end.SetString("999999999999999999999999999999", 10)
	fmt.Println(countZeroEndingCombinations(start, end))
}

#!/bin/bash

# 素数判定関数
is_prime() {
  local num=$1
  if [ "$num" -lt 2 ]; then
    return 1
  fi
  for ((i=2; i*i<=num; i++)); do
    if [ $((num % i)) -eq 0 ]; then
      return 1
    fi
  done
  return 0
}

# 指定範囲で素数を計算
start_range=1
end_range=50000
prime_count=0

echo "Finding primes between $start_range and $end_range..."
for ((num=start_range; num<=end_range; num++)); do
  if is_prime "$num"; then
    prime_count=$((prime_count + 1))
  fi
done

echo "Found $prime_count primes between $start_range and $end_range."

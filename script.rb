require 'benchmark'

# 素数判定メソッド
def prime?(n)
  return false if n < 2
  (2..Math.sqrt(n)).none? { |i| n % i == 0 }
end

# 指定範囲の素数を計算
def find_primes_in_range(start_num, end_num)
  primes = []
  (start_num..end_num).each do |num|
    primes << num if prime?(num)
  end
  primes
end

# 計算範囲を指定
start_range = 1
end_range = 100_000_0 # 範囲を大きくすると処理が重くなります

# 実行時間を測定
execution_time = Benchmark.measure do
  primes = find_primes_in_range(start_range, end_range)
  puts "Found #{primes.size} primes between #{start_range} and #{end_range}"
end

puts "Execution time: #{execution_time.real.round(2)} seconds"

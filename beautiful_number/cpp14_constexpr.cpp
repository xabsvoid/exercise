#include <cstdint>
#include <iostream>
#include <string>

// Возведение n в степень k
template <typename T>
constexpr auto exponentiation(T n, T k) -> T {
    auto e = T{1};
    for (T i = 0; i < k; i++) {
        e *= n;
    }
    return e;
}

// Размещение из множества n различных элементов по k элементов с повторениями
template <typename T>
constexpr auto permutation(T n, T k) -> T {
    return exponentiation<T>(n, k);
}

// Композиция натурального числа n на k слагаемых, каждое из которых не
// превосходит m
template <typename T>
constexpr auto composition(T n, T k, T m) -> T {
    if (n == 0 && k == 0) {
        return 1;
    }
    if (n > k * m) {
        return 0;
    }
    auto count = T{};
    for (T i = 0; i <= m; i++) {
        count += composition<T>(n - i, k - 1, m);
    }
    return count;
}

// Количество чисел сумма цифр левой и правой частей которого равна,
// где len - это длинна числа, num - система счисления
template <typename T>
constexpr auto beautiful_number(T len, T num) -> T {
    auto count = T{};
    auto k = len / 2;
    auto m = num - 1;
    auto max_sum = len * (num - 1);
    for (T sum = 0; sum <= max_sum; sum++) {
        count += permutation<T>(composition<T>(sum, k, m), 2);
    }
    if (len % 2 != 0) {
        count *= num;
    }
    return count;
}

constexpr auto beautiful_number_13_13 = beautiful_number<std::uint64_t>(13, 13);

int main(int argc, char *argv[]) {
    std::cout << std::to_string(beautiful_number_13_13) << std::endl;
    return 0;
}

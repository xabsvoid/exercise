#include <cstdint>
#include <iostream>
#include <string>

// Возведение n в степень k
template <typename T>
constexpr auto exponentiation(T n, T k) -> T {
    return k == 0 ? 1 : (n * exponentiation<T>(n, k - 1));
}

// Размещение из множества n различных элементов по k элементов с повторениями
template <typename T>
constexpr auto permutation(T n, T k) -> T {
    return exponentiation<T>(n, k);
}

// Композиция натурального числа n на k слагаемых, каждое из которых не
// превосходит m
template <typename T>
constexpr auto composition(T n, T k, T m, T mi = 0) -> T {
    return
        (n == 0 && k == 0) ? 1 :
        (n > k * m) ? 0 :
        (m == mi) ? composition<T>(n - mi, k - 1, m) :
                    composition<T>(n - mi, k - 1, m)
                    + composition<T>(n, k, m, mi + 1);
}

// Количество чисел сумма цифр левой и правой частей которого равна,
// где len - это длинна числа, num - система счисления
template <typename T>
constexpr auto beautiful_number(T len, T num, T sum = 0) -> T {
    return
        (permutation<T>(composition<T>(sum, len / 2, num - 1), 2)
        + ((sum < len * (num - 1)) ? beautiful_number(len, num, sum + 1) : 0))
        * ((sum == 0 && len % 2 != 0) ? num : 1);
}

constexpr auto beautiful_number_13_13 = beautiful_number<std::uint64_t>(13, 13);

int main(int argc, char *argv[]) {
    std::cout << std::to_string(beautiful_number_13_13) << std::endl;
    return 0;
}

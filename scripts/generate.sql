-- truncate all

truncate table user_solutions_results, user_solutions_tests, last_opened ,users, tests, tasks, user_solutions,
    tokens_for_password_reset, tokens_for_verification;

-- insert default user

insert into users (id, is_admin, first_name, last_name, username, email, password, activated)
values (0, false, '', '', '', '', '', false);

-- insert default task

insert into tasks (id, author_id, approver_id, title, difficulty, is_published, added_on,
                   text)
values (0, 0, 0, '', '', false, CURRENT_TIMESTAMP, '');

-- insert default test

insert into tests (id, last_modified, final, name, public, user_id, task_id, language, code)
values (0, CURRENT_TIMESTAMP, false, '', false, 0, 0, '', '');

-- insert default user solution

insert into user_solutions
    (id, user_id, task_id, last_modified, language, name, public, code)
values (0, 0, 0, CURRENT_TIMESTAMP, '', '', false, '');

-- restart all

alter sequence users_id_seq restart with 1;
alter sequence tasks_id_seq restart with 1;
alter sequence user_solutions_id_seq restart with 1;
alter sequence tests_id_seq restart with 1;

-- insert test users (passwords are 1234)

insert into users (is_admin, first_name, last_name, username, email, password, activated)
values (true, 'website', 'admin', 'admin', 'jan@jankelemen.net',
        '$2a$10$dY6ifXE0GuutyZwE0OjL/OsRcNLrI6N2HiZpaf.vD8/nAU7txxIX2', true),
       (true, 'second', 'admin', 'admin2', 'janko.kelemen@gmail.com',
        '$2a$10$6yl63KhFSNK0ds3eSxY6CONCXIwSznYZzlQi2h560cx9rT1VYDS9.', true),
       (false, 'first', 'user', 'user1', '	my.bachelors.thesis@gmail.com',
        '$2a$10$ZcLkQLciNXCq50cOHMSX2ORf7DHd0rRVdn7XmGZZHC37kdYQEa.Xa', true);

-- task Fizzbuzz

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Fizz buzz', 'easy',
        true, CURRENT_TIMESTAMP,
        '<p>Fizz buzz is a group word game for children to teach them about division. ' ||
        'Players take turns to count incrementally, replacing any number divisible by three with the word "fizz", and any number divisible by five with the word "buzz". ' ||
        'Write a function named "TestFizzBuzz1_000_000", which will return an array with the first 1 000 000 elements of this sequence.</p>');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 1, 'go', 'package main

import "testing"

func TestFizzBuzz1_000_000(t *testing.T) {
  res := FizzBuzz1_000_000()

  if len(res) != 1_000_000 {
    t.Errorf("got an array with length %d, want an array with length 1 000 000", len(res))
  }

  tests := []struct {
    name      string
    expecting string
    got       string
  }{
    {name: "on index 14", expecting: "fizzbuzz", got: res[14]},
    {name: "on index 100_000", expecting: "100001", got: res[100_000]},
    {name: "on index 999_999", expecting: "buzz", got: res[999_999]},
  }
  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      if test.expecting != test.got {
        t.Errorf("%s got: %q, expecting: %q", test.name, test.got, test.expecting)
      }
    })
  }
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 1, 'javascript', 'let assert = require("assert")

describe("Final test", function () {
  let got = FizzBuzz1_000_000()

  it("Test length", function () {
    assert.equal(got.length, 1_000_000)
  })

  it("Test indexes", function () {
    assert.equal(got[14], "fizzbuzz")
    assert.equal(got[100_000], "100001")
    assert.equal(got[999_999], "buzz")
  })
})');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 1, CURRENT_TIMESTAMP, 'go', 'public solution', true, 'package main

func FizzBuzz1_000_000() []string {
  // write your solution here
  return nil
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 1, CURRENT_TIMESTAMP, 'javascript', 'public solution', true, 'function FizzBuzz1_000_000() {
  return []
}');

-- task Primes

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Get first 1000 primes', 'easy',
        true, CURRENT_TIMESTAMP,
        '<p>Rewrite already an existing solution in Go into Python</p>');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 2, 'go', 'package main

import "testing"

func TestPrimes(t *testing.T) {
  got := primes()

  if len(got) != 1000 {
    t.Errorf("got an array with length %d, want an array with length 1000", len(got))
  }

  tests := []struct {
    name      string
    expecting int
    got       int
  }{
    {name: "on index 45", expecting: 197, got: got[45]},
    {name: "on index 234", expecting: 1481, got: got[234]},
    {name: "on index 980", expecting: 7723, got: got[980]},
  }
  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      if test.expecting != test.got {
        t.Errorf("%s got: %d, expecting: %d", test.name, test.got, test.expecting)
      }
    })
  }
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 2, 'python', 'def test_primes():
  got = primes()
  assert len(got) == 1000
  assert got[45] == 197
  assert got[234] == 1481
  assert got[980] == 7723');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 2, CURRENT_TIMESTAMP, 'go', 'solution in Go', true, 'package main

import (
  "math"
)

func primes() (res []int) {
  isPrime := func(n int) bool {
    limit := int(math.Pow(float64(n), 0.5)) + 1
    for i := 2; i < limit; i++ {
      if n%i == 0 {
        return false
      }
    }
    return true
  }

  for i := 1; len(res) <= 999; i++ {
    if isPrime(i) {
      res = append(res, i)
    }
  }
  return
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 2, CURRENT_TIMESTAMP, 'python', 'public solution', true, 'def primes():
  # write your solution here
  return []');

-- task Sum of two

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Sum of two', 'easy', true, CURRENT_TIMESTAMP,
        '<p>Write a function named "Sum", which will take two arguments - the first one is an array of numbers and the second one is a number. ' ||
        'Return indexes of two numbers from the array which sum to the second argument. If such numbers don''t exist, return an empty array</p>');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 3, 'go', 'package main

import (
  "reflect"
  "testing"
)

func TestSum(t *testing.T) {
  testTable := []struct {
    inputArr, want []int
    sum            int
  }{
    {inputArr: []int{2, 7, 11, 15}, sum: 9, want: []int{0, 1}},
    {inputArr: []int{3, 2, 4}, sum: 6, want: []int{1, 2}},
    {inputArr: []int{2, 7, 11, 15}, sum: 90, want: []int{}},
  }

  for _, test := range testTable {
    got := Sum(test.inputArr, test.sum)
    if !reflect.DeepEqual(got, test.want) {
      t.Errorf("got: %v, but want: %v. Input was %v and %d",
        got, test.want, test.inputArr, test.sum)
    }
  }
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 3, 'cpp', '#include <vector>

using namespace std;

TEST_CASE("Basic tests") {
  REQUIRE(Sum({2, 7, 11, 15}, 9) == vector<int>{0, 1});
  REQUIRE(Sum({3, 2, 4}, 6) == vector<int>{1, 2});
  REQUIRE(Sum({2, 7, 11, 15}, 90) == vector<int>{});
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 3, CURRENT_TIMESTAMP, 'go', 'public solution', true, 'package main

func Sum(arr []int, sum int) []int {
  // write your solution here
  return nil
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 3, CURRENT_TIMESTAMP, 'cpp', 'public solution', true, '#include <vector>

using namespace std;

vector<int> Sum(vector<int> arr, int sum) {
  // write your solution here
  return {};
}');

-- task Reverse Integer

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Reverse integer', 'medium', true, CURRENT_TIMESTAMP,
        '<p>Write a function named "Reverse", which for given x will return reversed x. Return 0 if x is not in range [-2^31, 2^31 - 1]. ' ||
        'For example:</p><p>-123 -> -321</p><p>100 -> 1</p><p>3 147 483 647 -> 0</p>');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 4, 'cpp', 'TEST_CASE("Basic reverse") {
  REQUIRE(Reverse(-1235) == -5321);
  REQUIRE(Reverse(98000) == 89);
}

TEST_CASE("Bigger numbers") {
  REQUIRE(Reverse(-147483648) == -846384741);
  REQUIRE(Reverse(-2147483649) == 0);
  REQUIRE(Reverse(214748364) == 463847412);
  REQUIRE(Reverse(2147483648) == 0);
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 4, 'javascript', 'let assert = require("assert")

describe("Final test", function () {
  it("Basic reverse", function () {
    assert.equal(Reverse(-1235), -5321)
    assert.equal(Reverse(98000), 89)
  }),
  it("Bigger numbers", function () {
    assert.equal(Reverse(-147483648), -846384741)
    assert.equal(Reverse(-2147483649), 0)
    assert.equal(Reverse(214748364), 463847412)
    assert.equal(Reverse(2147483648), 0)
  })
})');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 4, CURRENT_TIMESTAMP, 'javascript', 'public solution', true, 'function Reverse(n) {
  // write your solution here
  return 0
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 4, CURRENT_TIMESTAMP, 'cpp', 'public solution', true, 'using namespace std;

int Reverse(long int n) {
  // write your solution here
  return 0;
}');

-- task Median of two arrays

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Median of two arrays', 'medium', true, CURRENT_TIMESTAMP,
        '<p>Write a function named "Median" which takes two sorted arrays and returns their median. For examaple:</p>' ||
        '<p>[1,3], [2] -> 2</p><p>[1,2], [3,4] -> 2.5</p>');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 5, 'javascript', 'let assert = require("assert")

let firstBigArray = Array.from(Array(10_000_000).keys())
let secondBigArray = [...Array(1_000_000).keys()].map(x => x + 13)

describe("Final test", function () {
  it("Basic cases", function () {
    assert.equal(Median([1, 3], [2]), 2)
    assert.equal(Median([1, 2], [3, 4]), 2.5)
  }),
  it("Bigger arrays", function () {
    assert.equal(Median(firstBigArray, secondBigArray), 4499999.5)
  })
})');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 5, 'python', 'firstBigArray = [i for i in range(0, 10_000_000)]
secondBigArray = [i for i in range(13, 1_000_000 + 13)]

def test_median():
  assert Median([1, 3], [2]) == 2
  assert Median([1, 2], [3, 4]) == 2.5
  assert Median(firstBigArray, secondBigArray), 4499999.5');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 5, CURRENT_TIMESTAMP, 'javascript', 'public solution', true, 'function Median(a, b) {
  // write your solution here
  return 0
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 5, CURRENT_TIMESTAMP, 'python', 'public solution', true, 'def Median(a, b):
  # write your solution here
  return 0');


-- task Rome

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Roman numbers', 'medium', true, CURRENT_TIMESTAMP,
        '<p>Write a function named "Rome" that converts Arabic numbers into Roman. For examaple:</p>' ||
        '<p>II -> 2</p><p>MCMXCIV -> 1994</p>');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 6, 'go', 'package main

import (
  "testing"
)

func TestRome(t *testing.T) {
  testTable := []struct {
    input string
    want  int
  }{
    {input: "III", want: 3},
    {input: "LVIII", want: 58},
    {input: "MCMXCIV", want: 1994},
  }

  for _, test := range testTable {
    got := Rome(test.input)
    if got != test.want {
      t.Errorf("got: %d, but want: %d. Input was %q",
        got, test.want, test.input)
    }
  }
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 6, 'python', 'def test_rome():
  assert Rome("III") == 3
  assert Rome("LVIII") == 58
  assert Rome("MCMXCIV") == 1994');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 6, 'cpp', 'TEST_CASE("final test") {
  REQUIRE(Rome("III") == 3);
  REQUIRE(Rome("LVIII") == 58);
  REQUIRE(Rome("MCMXCIV") == 1994);
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 6, 'javascript', 'let assert = require("assert")

describe("Final test", function () {
  it("Basic cases", function () {
    assert.equal(Rome("III"), 3)
    assert.equal(Rome("LVIII"), 58)
    assert.equal(Rome("MCMXCIV"), 1994)
  })
})');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 6, CURRENT_TIMESTAMP, 'go', 'public solution', true, 'package main

func Rome(romeNumber string) int {
  // write your solution here
  return 0
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 6, CURRENT_TIMESTAMP, 'cpp', 'public solution', true, '#include <string>

using namespace std;

int Rome(string romeNumber) {
  // write your solution here
  return 0;
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 6, CURRENT_TIMESTAMP, 'javascript', 'public solution', true, 'function Rome(romeNumber) {
  // write your solution here
  return 0
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 6, CURRENT_TIMESTAMP, 'python', 'public solution', true, 'def Rome(romeNumber):
  # write your solution here
  return 0');

-- task Valid Parentheses

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Valid Parentheses', 'hard', true, CURRENT_TIMESTAMP,
        '<p>Write a function named "Parentheses" that returns a number of parentheses that aren''t mismatched. For examaple:</p>' ||
        '<p>(() -> 2</p><p>)[{()()}]) -> 8</p>');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 7, 'go', 'package main

import (
  "testing"
)

func TestParentheses(t *testing.T) {
  testTable := []struct {
    input string
    want  int
  }{
    {input: "(()", want: 2},
    {input: ")[()()])", want: 6},
    {input: "{{{{[{(abc)}]}()", want: 10},
  }

  for _, test := range testTable {
    got := Parentheses(test.input)
    if got != test.want {
      t.Errorf("got: %d, but want: %d. Input was %q",
        got, test.want, test.input)
    }
  }
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 7, 'cpp', 'TEST_CASE("final test") {
  REQUIRE(Parentheses("(()") == 2);
  REQUIRE(Parentheses(")[()()])") == 6);
  REQUIRE(Parentheses("{{{{[{(abc)}]}()") == 10);
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 7, 'javascript', 'let assert = require("assert")

describe("Final test", function () {
  it("Basic cases", function () {
    assert.equal(Parentheses("(()"), 2)
    assert.equal(Parentheses(")[()()])"), 6)
    assert.equal(Parentheses("{{{{[{(abc)}]}()"), 10)
  })
})');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 7, 'python', 'def test_parentheses():
  assert Parentheses("(()") == 3
  assert Parentheses(")[()()])") == 6
  assert Parentheses("{{{{[{(abc)}]}()") == 10');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 7, CURRENT_TIMESTAMP, 'go', 'public solution', true, 'package main

func Parentheses(parentheses string) int {
  // write your solution here
  return 0
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 7, CURRENT_TIMESTAMP, 'cpp', 'public solution', true, '#include <string>

using namespace std;

int Parentheses(string parentheses) {
  // write your solution here
  return 0;
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 7, CURRENT_TIMESTAMP, 'javascript', 'public solution', true, 'function Parentheses(parentheses) {
  // write your solution here
  return 0
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 7, CURRENT_TIMESTAMP, 'python', 'public solution', true, 'def Parentheses(parentheses):
  # write your solution here
  return 0');

-- task Longest Substring

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Longest Substring', 'medium', true, CURRENT_TIMESTAMP,
        '<p>Write a function named "Longest" that returns the longest substring without repeating characters. For example:</p>' ||
        '<p>abcabcbb -> 3 (because "abc")</p>' ||
        '<p>"bbbbb" -> 1 (because "b")</p>');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 8, 'go', 'package main

import (
  "strings"
  "testing"
)

func TestLongest(t *testing.T) {
  longString := strings.Repeat("a", 500_000) + "abcdef" + strings.Repeat("a", 500_000)

  testTable := []struct {
    input string
    want  int
  }{
    {input: "abcabcbb", want: 2},
    {input: "aaa", want: 1},
    {input: longString, want: 6},
  }

  for _, test := range testTable {
    got := Longest(test.input)
    if got != test.want {
      if len(test.input) > 10 {
        test.input = "too long to show"
      }
      t.Errorf("got: %d, but want: %d. Input was %q",
        got, test.want, test.input)
    }
  }
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 8, 'cpp', '#include <string>

using namespace std;

string getA() {
  string res = "";
  for (int i = 0; i < 500000; i++) {
    res += "a";
  }
  return res;
}

TEST_CASE("final test") {
  string longString = getA() + "abcdef" + getA();

  REQUIRE(Longest("abcabcbb") == 2);
  REQUIRE(Longest("aaa") == 1);
  REQUIRE(Longest(longString) == 6);
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final', true, 1, 8, 'javascript', 'let assert = require("assert")

let longString = "a".repeat(500_000) + "abcdef" + "a".repeat(500_000)

describe("Final test", function () {
  it("Basic cases", function () {
    assert.equal(Longest("abcabcbb"), 2)
    assert.equal(Longest("aaa"), 1)
    assert.equal(Longest(longString), 6)
  })
})');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 8, CURRENT_TIMESTAMP, 'go', 'public solution', true, 'package main

func Longest(s string) int {
  // write your solution here
  return 0
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 8, CURRENT_TIMESTAMP, 'cpp', 'public solution', true, '#include <string>

using namespace std;

int Longest(string s) {
  // write your solution here
  return 0;
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 8, CURRENT_TIMESTAMP, 'javascript', 'public solution', true, 'function Longest(s) {
  // write your solution here
  return 0
}');
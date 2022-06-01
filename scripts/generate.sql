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
values (true, 'website', 'admin', 'admin', 'admin@website.com',
        '$2a$10$dY6ifXE0GuutyZwE0OjL/OsRcNLrI6N2HiZpaf.vD8/nAU7txxIX2', true),
       (false, 'first', 'user', 'user1', 'user1@email.com',
        '$2a$10$6yl63KhFSNK0ds3eSxY6CONCXIwSznYZzlQi2h560cx9rT1VYDS9.', true),
       (false, 'second', 'user', 'user2', 'user2@email.com',
        '$2a$10$ZcLkQLciNXCq50cOHMSX2ORf7DHd0rRVdn7XmGZZHC37kdYQEa.Xa', true);

-- task Fizzbuzz

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Fizz buzz', 'easy',
        true, CURRENT_TIMESTAMP,
        'Fizz buzz is a group word game for children to teach them about division. ' ||
        'Players take turns to count incrementally, replacing any number divisible by three with the word "fizz", and any number divisible by five with the word "buzz". ' ||
        'Write a function named "TestFizzBuzz1_000_000", which will return an array with the first 1 000 000 elements of this sequence.');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 1, 'go', 'package main

import "testing"

func TestFizzBuzz1_000_000(t *testing.T) {
	res := FizzBuzz1_000_000()

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

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 1, CURRENT_TIMESTAMP, 'go', 'public solution', true, 'func FizzBuzz1_000_000() []int {
	return nil
}');

-- insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
-- values (1, 1, CURRENT_TIMESTAMP, 'go', 'my solution 1', false, 'package main
--
-- import (
-- 	"strconv"
-- )
--
-- func FizzBuzz1_000_000() []string {
-- 	res := make([]string, 0, 1_000_000)
-- 	for i := 1; i <= 1_000_000; i++ {
-- 		if i%3 == 0 && i%5 == 0 {
-- 			res = append(res, "fizzbuzz")
-- 		} else if i%3 == 0 {
-- 			res = append(res, "fizz")
-- 		} else if i%5 == 0 {
-- 			res = append(res, "buzz")
-- 		} else {
-- 			res = append(res, strconv.Itoa(i))
-- 		}
-- 	}
-- 	return res
-- }');

-- task Primes

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Get first 1000 primes', 'easy',
        true, CURRENT_TIMESTAMP,
        'Rewrite already an existing solution in Go into Python');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 2, 'go', 'package main

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
}
');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 2, 'python', 'def test_primes():
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
values (1, 2, CURRENT_TIMESTAMP, 'python', 'public solution', true, 'def primes() {
	# write your solution here
	return []
}');

-- task Sum of two

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Sum of two', 'easy', true, CURRENT_TIMESTAMP,
        'Write a function named "Sum", which will take two arguments - the first one is an array of numbers and the second one is a number. ' ||
        'Return indexes of two numbers from the array which sum to the second argument. If such numbers don''t exist, return an empty array');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 3, 'go', 'package main

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
}
');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 3, CURRENT_TIMESTAMP, 'go', 'public solution', true, 'package main

func Sum(arr []int, sum int) []int {
	// write your solution here
	return nil
}');

-- task Reverse Integer

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Reverse integer', 'medium', true, CURRENT_TIMESTAMP,
        'Write a function named "Reverse", which for given x will return reversed x. Return 0 if x is not in range [-2^31, 2^31 - 1].<br>' ||
        'For example:<br>-123 -> -321<br>100 -> 1<br>3 147 483 647 -> 0');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 4, 'cpp', 'TEST_CASE("basic reverse") {
    REQUIRE(Reverse(-1235) == -5321);
    REQUIRE(Reverse(98000) == 89);
}

TEST_CASE("Testing overflow") {
	REQUIRE(Reverse(-2147483648) == -8463847412);
	REQUIRE(Reverse(-2147483649) == 0);
	REQUIRE(Reverse(2147483647) == 7463847412);
	REQUIRE(Reverse(2147483648) == 0);
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 4, 'javascript', 'let assert = require("assert")

describe("Final test", function () {
  it("Basic reverse", function () {
    assert.equal(Reverse(-1235), -5321)
    assert.equal(Reverse(98000), 89)
  }),
  it("Testing overflow", function () {
    assert.equal(Reverse(-2147483648), -8463847412)
    assert.equal(Reverse(-2147483649), 0)
    assert.equal(Reverse(2147483647), 7463847412)
    assert.equal(Reverse(2147483648), 0)
  })
})');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 4, CURRENT_TIMESTAMP, 'javascript', 'public solution', true, 'function Reverse(n) {
  // write your solution here
  return 0
}');

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 4, CURRENT_TIMESTAMP, 'cpp', 'public solution', true, 'int Reverse(int n) {
  // write your solution here
  return 0;
}');

-- task Median of two arrays

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Median of two arrays', 'medium', true, CURRENT_TIMESTAMP,
        'Write a function named "Median" which takes two sorted arrays and returns their median<br>For examaple:' ||
        '<br>[1,3], [2] -> 2<br>[1,2], [3,4] -> 2.5');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 5, 'javascript', 'let assert = require("assert")

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

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 5, CURRENT_TIMESTAMP, 'javascript', 'public solution', true, 'function Median(a, b) {
  // write your solution here
  return 0
}');


-- task Rome

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Roman numbers', 'medium', true, CURRENT_TIMESTAMP,
        'Write a function named "Rome" that converts Arabic numbers into Roman<br>For examaple:' ||
        '<br>II -> 2<br>MCMXCIV -> 1994');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 6, 'go', 'package main

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
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 6, 'python', 'def test_rome():
  assert Rome("III") == 3
  assert Rome("LVIII") == 58
  assert Rome("MCMXCIV") == 1994');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 6, 'cpp', 'TEST_CASE("final test") {
  REQUIRE(Rome("III") == 3);
  REQUIRE(Rome("LVIII") == 58);
  REQUIRE(Rome("MCMXCIV") == 1994);
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 6, 'javascript', 'let assert = require("assert")

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
values (1, 6, CURRENT_TIMESTAMP, 'cpp', 'public solution', true, 'int Rome(string romeNumber) {
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
        'Write a function named "Parentheses" that returns a number of parentheses that aren''t mismatched<br>For examaple:<br>' ||
        '(() -> 2<br>)[{()()}]) -> 8');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 7, 'go', 'func TestParentheses(t *testing.T) {
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
}
');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 7, 'cpp', 'TEST_CASE("final test") {
  REQUIRE(Parentheses("(()") == 2);
  REQUIRE(Parentheses(")[()()])") == 6);
  REQUIRE(Parentheses("{{{{[{(abc)}]}()") == 10);
}');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 7, 'javascript', 'let assert = require("assert")

describe("Final test", function () {
  it("Basic cases", function () {
    assert.equal(Parentheses("(()"), 2)
    assert.equal(Parentheses(")[()()])"), 6)
    assert.equal(Parentheses("{{{{[{(abc)}]}()"), 10)
  })
})');

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 7, 'python', 'def test_parentheses():
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
values (1, 7, CURRENT_TIMESTAMP, 'cpp', 'public solution', true, 'int Parentheses(string parentheses) {
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

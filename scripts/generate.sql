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

-- insert fizzbuzz into tasks

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Fizz buzz', 'easy',
        true, CURRENT_TIMESTAMP,
        'Fizz buzz is a group word game for children to teach them about division. Players take turns to count incrementally, replacing any number divisible by three with the word "fizz", and any number divisible by five with the word "buzz".');

-- insert primes into tasks

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Get first 1000 primes', 'easy',
        true, CURRENT_TIMESTAMP,
        'Rewrite already an existing solution in Go into Python');

-- insert "some of two" into tasks

insert into tasks (author_id, approver_id, title, difficulty, is_published, added_on, text)
values (1, 1, 'Some of two', 'easy', true, CURRENT_TIMESTAMP,
        'Write a function named "Sum", which will take two arguments - the first one is an array of numbers and the second one is a number. Return indexes of two numbers from the array which sum to the second argument. If such numbers don''t exist, return an empty array');

-- insert "some of two" into tests
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

-- insert fizzbuzz into tests

insert into tests (last_modified, final, name, public, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 'final test', true, 1, 1, 'go', 'package main

import "testing"

func TestFizzBuzz1000(t *testing.T) {
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

-- insert primes into tests

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

-- insert fizzbuzz into user_solutions

insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
values (1, 1, CURRENT_TIMESTAMP, 'go', 'my solution 1', false, 'package main

import (
	"strconv"
)

func FizzBuzz1_000_000() []string {
	res := make([]string, 0, 1_000_000)
	for i := 1; i <= 1_000_000; i++ {
		if i%3 == 0 && i%5 == 0 {
			res = append(res, "fizzbuzz")
		} else if i%3 == 0 {
			res = append(res, "fizz")
		} else if i%5 == 0 {
			res = append(res, "buzz")
		} else {
			res = append(res, strconv.Itoa(i))
		}
	}
	return res
}');

-- insert primes into user_solutions

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

-- insert into user_solutions (user_id, task_id, last_modified, language, name, public, code)
-- values (1, 2, CURRENT_TIMESTAMP, 'python', 'my python solution', false, 'def is_prime(n):
-- 	for i in range(2, int(n**1 / 2) + 1):
-- 		if n % i == 0:
-- 			return False
-- 	return True
--
-- def primes():
-- 	res = []
-- 	i = 1
-- 	while len(res) < 100:
-- 		if is_prime(i):
-- 			res.append(i)
-- 		i += 1
-- 	return res');
--

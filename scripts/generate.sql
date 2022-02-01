-- truncate all

truncate table users, tests, tasks, user_solutions;

-- insert default user

insert into users (id, is_admin, first_name, last_name, username, email, password)
values (0, false, '', '', '', '', '');

-- insert default task

insert into tasks (id, author_id, approver_id, title, difficulty, description, is_published, is_approved, added_on,
                   text)
values (0, 0, 0, '', '', '', false, false, CURRENT_TIMESTAMP, '');

-- insert default test

insert into tests (id, last_modified, final, user_id, task_id, language, code)
values (0, CURRENT_TIMESTAMP, false, 0, 0, '', '');

-- restart all

alter sequence users_id_seq restart with 1;
alter sequence tasks_id_seq restart with 1;
alter sequence user_solutions_id_seq restart with 1;
alter sequence tests_id_seq restart with 1;

-- insert test users (passwords are 1234)

insert into users (is_admin, first_name, last_name, username, email, password)
values (true, 'Bill', 'The admin', 'admin', 'admin@bill.com',
        '$2a$10$dY6ifXE0GuutyZwE0OjL/OsRcNLrI6N2HiZpaf.vD8/nAU7txxIX2'),
       (false, 'Taylor', 'The user', 'taylor', 'taylor@email.com',
        '$2a$10$6yl63KhFSNK0ds3eSxY6CONCXIwSznYZzlQi2h560cx9rT1VYDS9.'),
       (false, 'Riley', 'Goodman', 'riley', 'riley@goodman.com',
        '$2a$10$ZcLkQLciNXCq50cOHMSX2ORf7DHd0rRVdn7XmGZZHC37kdYQEa.Xa');

-- insert fizzbuzz into tasks

insert into tasks (author_id, approver_id, title, difficulty, description, is_published, is_approved, added_on, text)
values (1, 1, 'Fizz buzz', 'easy', 'Fizz buzz is a group word game for children to teach them about division',
        true, true, CURRENT_TIMESTAMP,
        'Fizz buzz is a group word game for children to teach them about division. Players take turns to count incrementally, replacing any number divisible by three with the word "fizz", and any number divisible by five with the word "buzz".');

-- insert primes into tasks

insert into tasks (author_id, approver_id, title, difficulty, description, is_published, is_approved, added_on, text)
values (1, 1, 'Get first 100 primes', 'easy', 'Rewrite already an existing solution in Go into Python',
        true, true, CURRENT_TIMESTAMP,
        'Rewrite already an existing solution in Go into Python');

-- insert fizzbuzz into tests

insert into tests (last_modified, final, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 1, 1, 'go', 'package main

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

insert into tests (last_modified, final, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 1, 2, 'go', 'package main

import "testing"

func TestPrimes(t *testing.T) {
	got := primes()

	tests := []struct {
		name      string
		expecting int
		got       int
	}{
		{name: "on index 1", expecting: 2, got: got[1]},
		{name: "on index 9", expecting: 23, got: got[9]},
		{name: "on index 89", expecting: 461, got: got[89]},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expecting != test.got {
				t.Errorf("%s got: %d, expecting: %d", test.name, test.got, test.expecting)
			}
		})
	}
}');

insert into tests (last_modified, final, user_id, task_id, language, code)
values (CURRENT_TIMESTAMP, true, 1, 2, 'python', 'def test_primes():
	got = primes()
	assert got[1] == 2
	assert got[9] == 23
	assert got[89] == 461');

-- insert fizzbuzz into user_solutions

insert into user_solutions (user_id, task_id, test_id, last_modified, language, code, exit_code, output,
                            compilation_time, real_time, kernel_time, user_time, max_ram_usage, binary_size)
values (1, 1, 1, CURRENT_TIMESTAMP, 'go', 'package main

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
}', 0, '', 0, 0, 0, 0, 0, 0);

-- insert primes into user_solutions

insert into user_solutions (user_id, task_id, test_id, last_modified, language, code, exit_code, output,
                            compilation_time, real_time, kernel_time, user_time, max_ram_usage, binary_size)
values (1, 2, 2, CURRENT_TIMESTAMP, 'go', 'package main

import "math"

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

	for i := 1; len(res) <= 100; i++ {
		if isPrime(i) {
			res = append(res, i)
		}
	}
	return
}', 0, '', 0, 0, 0, 0, 0, 0);

insert into user_solutions (user_id, task_id, test_id, last_modified, language, code, exit_code, output,
                            compilation_time, real_time, kernel_time, user_time, max_ram_usage, binary_size)
values (1, 2, 3, CURRENT_TIMESTAMP, 'python', 'def is_prime(n):
	for i in range(2, int(n**1 / 2) + 1):
		if n % i == 0:
			return False
	return True

def primes():
	res = []
	i = 1
	while len(res) < 100:
		if is_prime(i):
			res.append(i)
		i += 1
	return res', 0, '', 0, 0, 0, 0, 0, 0);


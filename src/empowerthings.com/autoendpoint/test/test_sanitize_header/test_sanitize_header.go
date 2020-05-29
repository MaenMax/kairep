package main 

import(
	"fmt"
	"strings"
    "empowerthings.com/autoendpoint/utils"
)


type TestCase struct {
    Input             string  //
    Required_Labels []string  //
    Output            string  // Expected Output
    Error             error   // Expected Error
}

func main() {

    var test_cases []TestCase = []TestCase {
        // No quote, Comma, Space
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},


        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        

        // // No quote, Semi-column, Space
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=; p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==; p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=; p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==; p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=; p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==; p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw=",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=; p256=5i4j5lk45jkq4j5lkrWTJAWrw=",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==; p256=5i4j5lk45jkq4j5lkrWTJAWrw=",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=; p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==; p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        //// No quote, Comma, No Space
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=,p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=,p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=,p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        
        //// No quote, Semi-column, No Space
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz;p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=;p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==;p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz;p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=;p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==;p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz;p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=;p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==;p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        



        // // Simple quote, Comma, Space
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz', p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=', p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==', p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz', p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=', p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==', p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz', p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=', p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==', p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        // // Simple quote, Semi-Column, Space
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz'; p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz='; p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=='; p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz'; p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz='; p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=='; p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz'; p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz='; p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=='; p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        // // Simple quote, Comma, No Space
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz',p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=',p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==',p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz',p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=',p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==',p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz',p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=',p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==',p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        // // Simple quote, Semi-Column, No Space
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz';p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=';p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==';p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz';p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=';p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==';p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz';p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=';p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==';p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},




        // // Double Quote, Comma, Space
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        // // Double Quote, Semi-Column, Space
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\"; p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\"; p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\"; p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\"; p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\"; p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\"; p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\"; p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\"; p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\"; p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        // // Double Quote, Comma, No Space
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\",p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\",p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\",p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\",p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\",p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\",p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\",p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\",p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\",p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        // // Double Quote, Semi-Column, No Space
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\";p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\";p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\";p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\";p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\";p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\";p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz\";p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=\";p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==\";p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},




        // // Mixed Simple quote/No Quote Type I, Comma, Space
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz', p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=', p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==', p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz', p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=', p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==', p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz', p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=', p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==', p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        // // Partial Simple quote Type 2, Comma, Space
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},


        //// Partial Simple quote Type 1, Comma, No Space
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz',p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=',p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==',p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz',p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=',p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==',p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz',p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=',p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==',p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        // // Partial Simple quote Type 2, Comma, No Space
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=,p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=,p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=,p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},


        // // Partial Simple quote Type 2, Comma, No Space
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=,p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=,p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=,p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        // // Mix of comma and semi columns
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw; jh=5/i4+j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw; jh=5_i4-j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw;jh=5/i4+j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw; jh=5_i4-j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==,p256=5i4j5lk45jkq4j5lkrWTJAWrw;jh=5/i4+j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw; jh=5_i4-j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz,p256=5i4j5lk45jkq4j5lkrWTJAWrw=; jh=5/i4+j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw; jh=5_i4-j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=; p256=5i4j5lk45jkq4j5lkrWTJAWrw=, jh=5/i4+j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw; jh=5_i4-j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==; p256=5i4j5lk45jkq4j5lkrWTJAWrw=,jh=5/i4+j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw; jh=5_i4-j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz;p256=5i4j5lk45jkq4j5lkrWTJAWrw==, jh=5/i4+j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw; jh=5_i4-j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=;p256=5i4j5lk45jkq4j5lkrWTJAWrw==,jh=5/i4+j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw; jh=5_i4-j5lk45jkq4j5lkrWTJAWrw",nil},


        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=\"45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},


        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=\"5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=\"5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=\"5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==\"",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==\", p256=\"5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},


        // // Single quote case
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh='45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256='5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256='5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256='5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256='5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256='5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256='5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256='5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256='5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256='5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256=5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256=5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256=5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256=5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256=5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256=5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256=5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256=5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256=5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256='5i4j5lk45jkq4j5lkrWTJAWrw'",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256='5i4j5lk45jkq4j5lkrWTJAWrw='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256='5i4j5lk45jkq4j5lkrWTJAWrw=='",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256='5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256='5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256='5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256='5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256='5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256='5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz', p256='5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=', p256='5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==', p256='5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"dh=45kjwrek_ljqkwljrflkawrf-lksz; p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil},



        // // Negative cases that should lead to a rejection
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz= p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz== p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz= p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz== p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz= p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz== p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},

        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz= p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz== p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz= p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz== p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz= p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz== p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz p256==5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz= p256==5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz== p256==5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz p256==5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz= p256==5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz== p256==5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz p256==5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz= p256==5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz== p256==5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},


        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz p256==5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz= p256==5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz== p256==5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz p256==5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz= p256==5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz== p256==5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz p256==5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz= p256==5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh==45kjwrek/ljqkwljrflkawrf+lksz== p256==5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lkszp256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lkszp256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lkszp256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},


        TestCase{"dh45kjwrek/=ljqkwljrflkawrf+lkszp256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrekljqkwljrflkawrf+lksz=p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lkszp256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh45kjwrek/ljqkwljrflkawrf+lksz=p2565i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh45kjwrek/ljqkwljrflkawrf+lkszp2565i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},


        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,, p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, ,p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,   ,p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw=,",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw=,,",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{",dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{",,dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrf,lkawrf+lksz= p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh,=45kjwrek/ljqkwljrflkawrf+lksz== p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256,=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j,5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrfl,,kawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j,,5lkrWTJAWrw==",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},


        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz,, p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, ,p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==,   ,p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw=,",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw=,,",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{",dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw=",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{",,dh=45kjwrek/ljqkwljrflkawrf+lksz, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrf,lkawrf+lksz= p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh,=45kjwrek/ljqkwljrflkawrf+lksz== p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256,=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j,5lk45jkq4j5lkrWTJAWrw==",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrfl,,kawrf+lksz==, p256=5i4j5lk45jkq4j5lkrWTJAWrw==",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz==, p256=5i4j5lk45jkq4j,,5lkrWTJAWrw==",nil,"",utils.ERR_MALFORMED_HEADER},

        TestCase{"dh=, p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=,",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{",dh=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{",",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},

        TestCase{"dh=, p256=5i4j5lk45jkq4j5lkrWTJAWrw",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=,",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{",dh=",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{",",nil,"",utils.ERR_MALFORMED_HEADER},

        TestCase{"=45kjwrek/ljqkwljrflkawrf+lksz=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},

        TestCase{"=45kjwrek/ljqkwljrflkawrf+lksz=",nil,"",utils.ERR_MALFORMED_HEADER},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=",[]string{"dh","p256"},"",utils.ERR_MISSING_REQUIRED_LABEL},
        TestCase{"p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MISSING_REQUIRED_LABEL},


        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=,",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256,",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},

        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256=,",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256",nil,"",utils.ERR_MALFORMED_HEADER},
        TestCase{"dh=45kjwrek/ljqkwljrflkawrf+lksz=, p256,",nil,"",utils.ERR_MALFORMED_HEADER},

        TestCase{"=45kjwrek/ljqkwljrflkawrf+lksz=, p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},

        TestCase{"45kjwrekljqkwljrflkawrflksz=p256=5i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"45kjwrekljqkwljrflkawrflksz=p2565i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MALFORMED_HEADER},
        TestCase{"f45kjwre=p2565i4j5lk45jkq4j5lkrWTJAWrw",[]string{"dh","p256"},"",utils.ERR_MISSING_REQUIRED_LABEL},

        TestCase{"f45kjwre=p2565i4j5lk45jkq4j5lkrWTJAWrw",nil,"f45kjwre=p2565i4j5lk45jkq4j5lkrWTJAWrw",nil},
        TestCase{"f45kjwre=p2565i4j5lk45jkq4j5lkrWTJAWrw",[]string{},"f45kjwre=p2565i4j5lk45jkq4j5lkrWTJAWrw",nil},

        TestCase{"dh=BKXED7dmK2xd202NNYGqlSLgTaBxNYOji8GFCjJ0XXnuqODbP9ma5bNuN6c0UL3vchL6wzJOxEKW_a48HWgGqG0=",nil,"dh=BKXED7dmK2xd202NNYGqlSLgTaBxNYOji8GFCjJ0XXnuqODbP9ma5bNuN6c0UL3vchL6wzJOxEKW_a48HWgGqG0",nil},

        TestCase{"keyid=p256dh;dh=BJmM2h717YDhOstWrLUi31mrhrcR8iKKE2-1aX1ZKP1zYp_lJZQs4O3LH5CpoKOjL66J0Ir25EU0oAWb-2yI6ns==;p256ecdsa=BGeoNWH-aSFuu-Bq5DCy1JGF021awFq5QrsHKMwC45wRtexY4JLYUVS9eb3JxswaSEeA8az1HzcsLrmuNuzM-uc",nil,"keyid=p256dh; dh=BJmM2h717YDhOstWrLUi31mrhrcR8iKKE2-1aX1ZKP1zYp_lJZQs4O3LH5CpoKOjL66J0Ir25EU0oAWb-2yI6ns; p256ecdsa=BGeoNWH-aSFuu-Bq5DCy1JGF021awFq5QrsHKMwC45wRtexY4JLYUVS9eb3JxswaSEeA8az1HzcsLrmuNuzM-uc",nil},

        }
    

    err_count:=0
    total_count:=0

    for i,tc:= range test_cases {
        total_count++
        // if i>0 {
        //     break
        // }

        fmt.Printf("%3d '%s' ", i, tc.Input)

        output,err:=utils.Sanitize_Header(tc.Input,tc.Required_Labels)

        if err!=tc.Error {
            err_count++
            fmt.Printf("KO (Wrong Returned Error '%v' instead of '%v')\n",err, tc.Error)
            continue
        }
        
        if strings.Compare(tc.Output,output)!=0 {
            err_count++
            fmt.Printf("KO (Wrong returned Output '%s' instead of '%s')\n",output,tc.Output)
        }

        fmt.Printf("OK\n")
    }

    fmt.Printf("------------------------\n")
    fmt.Printf("Total: %d, Error: %d\n",total_count, err_count)
}



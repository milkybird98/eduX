echo ">>>>1000 concurrency ,2 request depth ,100k request total."
time go test par_a_test.go
echo ""
echo ">>>>50 concurrency ,5 request depth ,100k request total."
time go test par_b_test.go

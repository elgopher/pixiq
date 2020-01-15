# How to write good automated tests

* first of all unit tests are the king. The ratio of unit tests should be at least 95% of the total number of tests. No kidding, that's the goal.
* tests should run fast - minute or two at most when we reach 1.0.0. If we have mostly unit tests then it should not be a problem
* for integration testing in CI we can use things like - software drivers using solely CPU (Mesa 3D), fake display servers such as XVFB


# Manual testing

* unfortunetly it is not possible to have only automated tests.
* some tests require manual work. 
* before merging each PR it should be tested manually using all [examples](examples)

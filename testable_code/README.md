Testable Codes in Golang
- Writing tests are not easy, writing comprehensive tests even harder, covering the critical paths then continue with non-critical ones would provide enough coverage for a reasonably robust code
- If code is easy to test then it is easy to understand, if it is hard to understand then it is hard to test, and vice versa
- Code coverage that covers the lines of code provides the minimum indicator that the lines have been run on test codes but does not tell if minimum or all conditions have been tested, hence 100% coverage is not equal to 100% tested, in fact, it is possible that 80% code coverage with more tests can actually covers more test cases than 100% coverage

The following are the guidelines/findings to reach the goal:
1. When defining a function instead of a method, be sure that it does need to be a function since a function nested inside another method or function cannot be stubbed
2. When defining a function / method, there are three types of parameters:
  - literal values
  - pure struct: a pure struct has not method, hence, only its fields need to be changed during testing
  - struct with methods: this struct has methods, hence, always define an interface and use this interface instead so that the struct methods can be stubbed during testing
3a. any variable that you want to test inside function/method will need to be returned so that the its value can be tested, however, there are variables that are 'in the middle' that is produced by some logic that should be tested - this logic needs to be extracted out as a method so that it can be tested in isolation
3b. If there a few possible values to test, put the possible inputs in a collection to iterate so that it can be done efficiently in one test method
4. Whenever you define a nested method, this method needs to be in an interface. One of the main challenges are grouping these nested methods in a meaningful interface definition. The art of coding is to balance between having multiple nested methods or can the code be made flatter - less level while still have the code testable
5. Testing conditional logic where one case does not have a value to test for, e.g: logger.Log("this case is skipped"). In go, reflection cannot help the issue either, so the green path can be let go without being tested - we don't know if a success path actually called the Log method
6. When possible, make use of custom error implementation so that you can test a specific error during testing instead of generic one, which can cause 'imposter' error
7.  OS/Db operations, external service calls need to be wrapped as well so that they can stubbed too
8. Checking error is very important in go as there is no error propagation, error raised can easily go silent
9. The nested method needs to be called explicity called with its interface so that the stubbed method can be used
10. In order to stub, the actual method needs to call its method using explicit interface, e.g: w.Worker.Run()
11. For naming test method, it can be named by the method being tested followed by any nested method then followed by description of what state causes the method being run. This usually is one of the if conditions, followed by response status if any, e.g: Test_Run_CheckStatus_CategoryA_Success() . If you can find the line / logic being tested easily by reading the test method name then it is named well enough
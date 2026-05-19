# Frequently Asked Questions:

### Q: Why did visitor service tests fail with nil pointer dereferences?
**A:** The `visitors.v1` Service depends on `devices.v1` and `visitor_sessions.v1` repositories to check and create devices and sessions. In unit tests, passing `nil` values for these repositories caused nil pointer dereferences when methods like `CreateSession` executed tracking logic.

### Q: How do we resolve dependency issues in unit tests?
**A:** We define mock repositories using Mockery. Add the relevant interfaces to `.mockery.yml` under the respective module names, run `mockery` to regenerate mocks, and instantiate/mock their expectations inside the unit tests.

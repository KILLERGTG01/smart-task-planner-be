# Contributing to Smart Task Planner Backend

## Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/smart-task-planner-be.git`
3. Install dependencies: `make install-deps`
4. Copy environment: `cp .env.example .env`
5. Run migrations: `make migrate`
6. Start development: `make dev`

## Code Standards

- Follow Go best practices
- Write tests for new features
- Use structured logging with Zap
- Validate all inputs
- Handle errors gracefully

## Pull Request Process

1. Create a feature branch
2. Make your changes
3. Add tests if applicable
4. Ensure all tests pass
5. Submit a pull request

## Reporting Issues

Please use GitHub Issues to report bugs or request features.
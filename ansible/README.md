# Docker Compose Service Tests with Ansible

This directory contains Ansible playbooks and roles to test the Docker Compose services defined in the root `docker-compose.yml` file.

## Prerequisites

- Ansible 2.9+
- Python 3.6+
- Docker and Docker Compose installed on the host
- Required Python packages: `docker`, `docker-compose`

## Directory Structure

```
ansible/
├── group_vars/
│   └── all.yml           # Global variables and service configurations
├── roles/
│   └── docker_services/
│       ├── tasks/
│       │   ├── main.yml              # Main tasks
│       │   ├── parse_compose.yml      # Parse Docker Compose file
│       │   ├── verify_services.yml    # Verify services are running
│       │   ├── verify_service.yml     # Verify individual service
│       │   ├── verify_dependencies.yml # Check service dependencies
│       │   ├── test_health_checks.yml # Test service health endpoints
│       │   └── test_connectivity.yml  # Test network connectivity
├── playbook.yml          # Main playbook
└── README.md             # This file
```

## Usage

1. Install required Python packages:
   ```bash
   pip install -r requirements.txt
   ```

2. Run all tests:
   ```bash
   ansible-playbook playbook.yml -i localhost,
   ```

3. Run specific test tags:
   ```bash
   # Test only Docker Compose setup
   ansible-playbook playbook.yml -i localhost, --tags docker_compose
   
   # Test service status
   ansible-playbook playbook.yml -i localhost, --tags services
   
   # Test health checks
   ansible-playbook playbook.yml -i localhost, --tags health_checks
   
   # Test network connectivity
   ansible-playbook playbook.yml -i localhost, --tags connectivity
   ```

## Customization

You can customize the test configuration by modifying the following files:

- `group_vars/all.yml`: Modify service ports, expected services, and health check endpoints
- `roles/docker_services/tasks/*.yml`: Customize test cases and validations

## Adding New Services

To add a new service to the test suite:

1. Add the service configuration to `group_vars/all.yml` under `expected_services`
2. Add health check endpoints under `health_checks` if applicable
3. Update any dependencies in the service configuration

## Troubleshooting

- If you encounter connection issues, ensure Docker is running and the services are up
- Check the Ansible output for detailed error messages
- Run with `-v` or `-vvv` for more verbose output

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

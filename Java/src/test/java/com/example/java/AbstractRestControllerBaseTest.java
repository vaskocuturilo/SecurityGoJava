package com.example.java;

import com.example.java.config.UserAuthenticationProvider;
import com.example.java.dto.UserDto;
import com.example.java.entity.UserEntity;
import com.example.java.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.postgresql.PostgreSQLContainer;

import java.util.UUID;

public abstract class AbstractRestControllerBaseTest {

    @Autowired
    private UserAuthenticationProvider userAuthenticationProvider;

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private PasswordEncoder passwordEncoder;

    private static final String TEST_EMAIL = "title@title.com";

    private static final String TEST_PASSWORD = "description";

    @BeforeEach
    void setupTestUser() {
        if (userRepository.findByEmail(TEST_EMAIL).isEmpty()) {
            final UserEntity user = new UserEntity();
            user.setEmail(TEST_EMAIL);
            user.setUsername(TEST_EMAIL);
            user.setPassword(passwordEncoder.encode(TEST_PASSWORD));
            user.setEnabled(true);
            userRepository.save(user);
        }
    }

    protected String generateTestToken() {
        final UserDto testUser = UserDto.builder()
                .id(UUID.randomUUID().toString())
                .email(TEST_EMAIL)
                .username(TEST_PASSWORD)
                .active(true)
                .build();
        return userAuthenticationProvider.createToken(testUser);
    }

    @Container
    static final PostgreSQLContainer POSTGRES_SQL_CONTAINER;

    static {
        POSTGRES_SQL_CONTAINER = new PostgreSQLContainer("postgres:17")
                .withUsername("postgres")
                .withPassword("password")
                .withDatabaseName("tasks_testcontainers");

        POSTGRES_SQL_CONTAINER.start();

    }

    @DynamicPropertySource
    public static void dynamicPropertySource(final DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", POSTGRES_SQL_CONTAINER::getJdbcUrl);
        registry.add("spring.datasource.username", POSTGRES_SQL_CONTAINER::getUsername);
        registry.add("spring.datasource.password", POSTGRES_SQL_CONTAINER::getPassword);
    }
}

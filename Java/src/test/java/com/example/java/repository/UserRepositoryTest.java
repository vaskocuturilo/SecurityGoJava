package com.example.java.repository;

import com.example.java.entity.UserEntity;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.data.jpa.test.autoconfigure.DataJpaTest;
import org.springframework.boot.jdbc.test.autoconfigure.AutoConfigureTestDatabase;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.test.context.junit.jupiter.SpringExtension;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;

@ExtendWith(SpringExtension.class)
@DataJpaTest(properties = {
        "spring.jpa.properties.javax.persistence.validation.mode=none",
        "spring.jpa.hibernate.ddl-auto=create-drop"
})
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
@Testcontainers
class UserRepositoryTest {

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:17");

    @Autowired
    private UserRepository userRepository;

    private static final String TEST_EMAIL = "title@title.com";

    private static final String TEST_PASSWORD = "description";

    @BeforeEach
    void setupTestUser() {
        userRepository.deleteAll();

        final UserEntity user = new UserEntity();
        user.setEmail(TEST_EMAIL);
        user.setUsername(TEST_EMAIL);
        user.setPassword(TEST_PASSWORD);
        user.setEnabled(true);
        userRepository.save(user);
    }

    @Test
    @DisplayName("Test findByEmail user functionality")
    void givenUserCreated_whenFindByEmail_thenUserIsReturned() {
        //given

        //when
        final UserEntity existUser = userRepository.findByEmail(TEST_EMAIL).orElse(null);

        //then
        assertThat(existUser).isNotNull();
        assertThat(existUser.getEmail()).isEqualTo(TEST_EMAIL);
    }

    @Test
    @DisplayName("Test findByEmail returns empty when user does not exist")
    void givenNoUser_whenFindByEmail_thenEmptyIsReturned() {
        //given

        //when
        final Optional<UserEntity> result = userRepository.findByEmail("nobody@test.com");

        //then
        assertThat(result).isEmpty();
    }
}
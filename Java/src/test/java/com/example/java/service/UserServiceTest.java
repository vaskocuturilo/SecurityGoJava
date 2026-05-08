package com.example.java.service;

import com.example.java.dto.CredentialsDto;
import com.example.java.dto.SignUpDto;
import com.example.java.dto.UserDto;
import com.example.java.entity.UserEntity;
import com.example.java.exception.UserException;
import com.example.java.mapper.UserMapper;
import com.example.java.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.BDDMockito;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.nio.CharBuffer;
import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;

@ExtendWith(MockitoExtension.class)
class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private UserService userService;

    @Mock
    private UserMapper userMapper;

    @Mock
    private PasswordEncoder passwordEncoder;

    private UserEntity testUser;

    private UserDto testUserDto;

    private static final String TEST_EMAIL = "title@title.com";

    private static final String TEST_PASSWORD = "description";

    @BeforeEach
    void setupTestUser() {
        testUser = new UserEntity();
        testUser.setId(UUID.randomUUID().toString());
        testUser.setEmail(TEST_EMAIL);
        testUser.setUsername(TEST_EMAIL);
        testUser.setPassword(TEST_PASSWORD);
        testUser.setEnabled(true);

        testUserDto = new UserDto(UUID.randomUUID().toString(), TEST_EMAIL, TEST_EMAIL, true);
    }

    @Test
    @DisplayName("Test the login user functionality")
    void givenValidCredential_whenLogin_thenReturnUserDto() {
        //given
        CredentialsDto credentials = new CredentialsDto(TEST_EMAIL, TEST_PASSWORD.toCharArray());

        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.of(testUser));
        BDDMockito.given(passwordEncoder.matches(any(CharBuffer.class), eq(TEST_PASSWORD))).willReturn(true);
        BDDMockito.given(userMapper.toUserDto(testUser)).willReturn(testUserDto);

        //when
        final UserDto result = userService.login(credentials);

        //then
        assertThat(result).isNotNull();
        assertThat(result.email()).isEqualTo(TEST_EMAIL);
        verify(userRepository).findByEmail(TEST_EMAIL);
        verify(passwordEncoder).matches(any(CharBuffer.class), eq(TEST_PASSWORD));
    }

    @Test
    @DisplayName("Test the login with unknown email functionality")
    void givenUnknownEmail_whenLogin_thenReturnException() {
        //given
        CredentialsDto credentials = new CredentialsDto(TEST_EMAIL, TEST_PASSWORD.toCharArray());

        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.empty());

        //when
        final UserException exception = assertThrows(UserException.class, () -> userService.login(credentials));

        //then
        assertThat(exception.getStatus()).isEqualTo(HttpStatus.NOT_FOUND);
        assertThat(exception.getMessage()).isEqualTo("Unknown user");
        verify(passwordEncoder, never()).matches(any(), any());
    }

    @Test
    @DisplayName("Test the login with disabled user functionality")
    void givenDisabledUser_whenLogin_thenReturnForbidden() {
        // given
        testUser.setEnabled(false);
        final CredentialsDto credentials = new CredentialsDto(TEST_EMAIL, "Wrong password".toCharArray());

        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.of(testUser));

        // when
        final UserException exception = assertThrows(UserException.class,
                () -> userService.login(credentials));

        // then
        assertThat(exception.getStatus()).isEqualTo(HttpStatus.FORBIDDEN);
        assertThat(exception.getMessage()).isEqualTo("The user is not enabled");
        verify(passwordEncoder, never()).matches(any(), any());
    }

    @Test
    @DisplayName("Test the login with wrong password functionality")
    void givenWrongPassword_whenLogin_thenReturnBadRequest() {
        // given
        final CredentialsDto credentials = new CredentialsDto(TEST_EMAIL, TEST_PASSWORD.toCharArray());

        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.of(testUser));

        BDDMockito.given(passwordEncoder.matches(any(CharBuffer.class), eq(TEST_PASSWORD))).willReturn(false);

        // when
        final UserException exception = assertThrows(UserException.class, () -> userService.login(credentials));

        // then
        assertThat(exception.getStatus()).isEqualTo(HttpStatus.BAD_REQUEST);
        assertThat(exception.getMessage()).isEqualTo("Invalid password");
    }


    @Test
    @DisplayName("Test the register new user functionality")
    void givenNewUser_whenRegister_thenSaveAndReturnUserDto() {
        // given
        final SignUpDto signUpDto = new SignUpDto(TEST_EMAIL, TEST_PASSWORD.toCharArray());

        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.empty());

        BDDMockito.given(passwordEncoder.encode(any(CharBuffer.class))).willReturn(TEST_PASSWORD);

        BDDMockito.given(userRepository.save(any(UserEntity.class))).willReturn(testUser);

        BDDMockito.given(userMapper.toUserDto(testUser)).willReturn(testUserDto);

        // when
        final UserDto result = userService.register(signUpDto);

        // then
        assertThat(result).isNotNull();
        verify(userRepository).save(argThat(entity -> {
            if (!entity.getEmail().equals(TEST_EMAIL) ||
                    !entity.getUsername().equals(TEST_EMAIL)) return false;
            assert entity.getPassword() != null;
            return entity.getPassword().equals(TEST_PASSWORD);
        }));
    }

    @Test
    @DisplayName("Test the register with duplicate email functionality")
    void givenExistingEmail_whenRegister_thenReturnBadRequest() {
        // given
        final SignUpDto signUpDto = new SignUpDto(TEST_EMAIL, TEST_PASSWORD.toCharArray());

        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.of(testUser));

        // when
        final UserException exception = assertThrows(UserException.class, () -> userService.register(signUpDto));

        // then
        assertThat(exception.getStatus()).isEqualTo(HttpStatus.BAD_REQUEST);
        assertThat(exception.getMessage()).isEqualTo("Email already exists");
        verify(userRepository, never()).save(any());
        verify(passwordEncoder, never()).encode(any());
    }

    @Test
    @DisplayName("Test the findByLogin functionality")
    void givenExistingEmail_whenFindByLogin_thenReturnUserDto() {
        // given
        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.of(testUser));

        BDDMockito.given(userMapper.toUserDto(testUser)).willReturn(testUserDto);

        // when
        final UserDto result = userService.findByLogin(TEST_EMAIL);

        // then
        assertThat(result).isNotNull();
        assertThat(result.email()).isEqualTo(TEST_EMAIL);
    }

    @Test
    @DisplayName("Test the findByLogin with unknown email functionality")
    void givenUnknownEmail_whenFindByLogin_thenReturnNotFound() {
        // given
        BDDMockito.given(userRepository.findByEmail("nobody@test.com")).willReturn(Optional.empty());

        // when / then
        final UserException exception = assertThrows(UserException.class, () -> userService.findByLogin("nobody@test.com"));

        assertThat(exception.getStatus()).isEqualTo(HttpStatus.NOT_FOUND);
        assertThat(exception.getMessage()).isEqualTo("Unknown user");
    }
}
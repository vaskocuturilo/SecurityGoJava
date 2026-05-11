package com.example.java.service;

import com.example.java.dto.RefreshTokenRequest;
import com.example.java.entity.RefreshTokenEntity;
import com.example.java.entity.UserEntity;
import com.example.java.exception.UserException;
import com.example.java.repository.RefreshTokenRepository;
import com.example.java.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.BDDMockito;
import org.mockito.InOrder;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.http.HttpStatus;

import java.time.LocalDateTime;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.AssertionsForClassTypes.assertThatNoException;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class RefreshTokenServiceTest {

    @Mock
    private RefreshTokenRepository refreshTokenRepository;

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private RefreshTokenService refreshTokenService;

    private UserEntity testUser;
    private RefreshTokenEntity testToken;

    private static final String TEST_EMAIL = "title@title.com";
    private static final String TEST_TOKEN = "test-refresh-token";

    @BeforeEach
    void setUp() {
        testUser = new UserEntity();
        testUser.setId("test-id");
        testUser.setEmail(TEST_EMAIL);
        testUser.setUsername(TEST_EMAIL);
        testUser.setPassword("hashed-password");
        testUser.setEnabled(true);

        testToken = RefreshTokenEntity.builder()
                .token(TEST_TOKEN)
                .user(testUser)
                .expiresAt(LocalDateTime.now().plusDays(7))
                .revoked(false)
                .build();
    }

    @Test
    @DisplayName("Test create refresh token")
    void givenValidEmail_whenCreateRefreshToken_thenTokenIsReturned() {
        // given
        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.of(testUser));
        BDDMockito.given(refreshTokenRepository.save(any(RefreshTokenEntity.class))).willReturn(testToken);

        // when
        final RefreshTokenEntity result = refreshTokenService.createRefreshToken(TEST_EMAIL);

        // then
        assertThat(result).isNotNull();
        assertThat(result.getUser()).isEqualTo(testUser);
        assertThat(result.isRevoked()).isFalse();
        assertThat(result.getExpiresAt()).isAfter(LocalDateTime.now());
        verify(refreshTokenRepository).deleteByUser(testUser);
        verify(refreshTokenRepository).save(any(RefreshTokenEntity.class));
    }

    @Test
    @DisplayName("Test create refresh token with unknow email")
    void givenUnknownEmail_whenCreateRefreshToken_thenThrowUserExceptionNotFound() {
        // given
        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.empty());

        // when
        final UserException exception = assertThrows(UserException.class, () -> refreshTokenService.createRefreshToken(TEST_EMAIL));

        //then
        assertThat(exception.getStatus()).isEqualTo(HttpStatus.NOT_FOUND);
        assertThat(exception.getMessage()).isEqualTo("Unknow user");
        verify(refreshTokenRepository, never()).deleteByUser(any());
        verify(refreshTokenRepository, never()).save(any());
    }

    @Test
    @DisplayName("Test deleted refresh token")
    void givenExistingToken_whenCreateRefreshToken_thenOldTokenDeletedFirst() {
        // given
        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.of(testUser));
        BDDMockito.given(refreshTokenRepository.save(any(RefreshTokenEntity.class))).willReturn(testToken);

        // when
        refreshTokenService.createRefreshToken(TEST_EMAIL);

        // then
        final InOrder inOrder = BDDMockito.inOrder(refreshTokenRepository);
        inOrder.verify(refreshTokenRepository).deleteByUser(testUser);
        inOrder.verify(refreshTokenRepository).save(any(RefreshTokenEntity.class));
    }

    @Test
    @DisplayName("Test create refresh token with valid UUID")
    void givenValidEmail_whenCreateRefreshToken_thenTokenIsUUID() {
        // given
        BDDMockito.given(userRepository.findByEmail(TEST_EMAIL)).willReturn(Optional.of(testUser));
        BDDMockito.given(refreshTokenRepository.save(any(RefreshTokenEntity.class))).willAnswer(invocation -> invocation.getArgument(0));

        // when
        final RefreshTokenEntity result = refreshTokenService.createRefreshToken(TEST_EMAIL);

        // then
        assertThatNoException().isThrownBy(() -> UUID.fromString(result.getToken()));
    }

    @Test
    @DisplayName("Test verify expiration")
    void givenValidToken_whenVerifyExpiration_thenTokenIsReturned() {
        // given
        BDDMockito.given(refreshTokenRepository.findByToken(TEST_TOKEN)).willReturn(Optional.of(testToken));

        // when
        final RefreshTokenEntity result = refreshTokenService.verifyExpiration(TEST_TOKEN);

        // then
        assertThat(result).isEqualTo(testToken);
        verify(refreshTokenRepository, never()).delete(any());
    }

    @Test
    @DisplayName("Test verify expiration with unknown token")
    void givenUnknownToken_whenVerifyExpiration_thenThrowUserExceptionUnauthorized() {
        // given
        BDDMockito.given(refreshTokenRepository.findByToken(TEST_TOKEN)).willReturn(Optional.empty());

        // when
        final UserException exception = assertThrows(UserException.class, () -> refreshTokenService.verifyExpiration(TEST_TOKEN));

        //then
        assertThat(exception.getStatus()).isEqualTo(HttpStatus.UNAUTHORIZED);
        assertThat(exception.getMessage()).isEqualTo("Invalid refresh token");
    }

    @Test
    @DisplayName("Test verify expiration revoked")
    void givenRevokedToken_whenVerifyExpiration_thenThrowUserExceptionUnauthorized() {
        // given
        testToken.setRevoked(true);

        BDDMockito.given(refreshTokenRepository.findByToken(TEST_TOKEN))
                .willReturn(Optional.of(testToken));

        // when
        final UserException exception = assertThrows(UserException.class, () -> refreshTokenService.verifyExpiration(TEST_TOKEN));

        //then
        assertThat(exception.getStatus()).isEqualTo(HttpStatus.UNAUTHORIZED);
        assertThat(exception.getMessage()).isEqualTo("Refresh token was revoked");
        verify(refreshTokenRepository, never()).delete(any());
    }

    @Test
    @DisplayName("Test verify expiration expired token")
    void givenExpiredToken_whenVerifyExpiration_thenDeleteAndThrowUserExceptionUnauthorized() {
        // given
        testToken.setExpiresAt(LocalDateTime.now().minusDays(1));
        BDDMockito.given(refreshTokenRepository.findByToken(TEST_TOKEN)).willReturn(Optional.of(testToken));

        // when
        final UserException exception = assertThrows(UserException.class, () -> refreshTokenService.verifyExpiration(TEST_TOKEN));

        //then
        assertThat(exception.getStatus()).isEqualTo(HttpStatus.UNAUTHORIZED);
        assertThat(exception.getMessage()).isEqualTo("Refresh token has expired, please log in again");
        verify(refreshTokenRepository).delete(testToken);
    }

    @Test
    @DisplayName("Test logout")
    void givenValidToken_whenLogout_thenTokenIsRevoked() {
        // given
        final RefreshTokenRequest request = new RefreshTokenRequest(TEST_TOKEN);

        BDDMockito.given(refreshTokenRepository.findByToken(TEST_TOKEN))
                .willReturn(Optional.of(testToken));

        // when
        refreshTokenService.logout(request);

        // then
        assertThat(testToken.isRevoked()).isTrue();
        verify(refreshTokenRepository).save(testToken);
    }

    @Test
    @DisplayName("Test logout with unknown token")
    void givenUnknownToken_whenLogout_thenThrowUserExceptionUnauthorized() {
        // given
        final RefreshTokenRequest request = new RefreshTokenRequest(TEST_TOKEN);

        BDDMockito.given(refreshTokenRepository.findByToken(TEST_TOKEN)).willReturn(Optional.empty());

        // when
        final UserException exception = assertThrows(UserException.class, () -> refreshTokenService.logout(request));

        //then
        assertThat(exception.getStatus()).isEqualTo(HttpStatus.UNAUTHORIZED);
        assertThat(exception.getMessage()).isEqualTo("Invalid refresh token");
        verify(refreshTokenRepository, never()).save(any());
    }

    @Test
    @DisplayName("Test revoke by user")
    void givenUserWithTokens_whenRevokeByUser_thenAllTokensRevoked() {
        // given
        final RefreshTokenEntity secondToken = RefreshTokenEntity.builder()
                .token("second-token")
                .user(testUser)
                .expiresAt(LocalDateTime.now().plusDays(3))
                .revoked(false)
                .build();

        BDDMockito.given(refreshTokenRepository.findAllByUser(testUser)).willReturn(List.of(testToken, secondToken));

        // when
        refreshTokenService.revokeByUser(testUser);

        // then
        assertThat(testToken.isRevoked()).isTrue();
        assertThat(secondToken.isRevoked()).isTrue();
        verify(refreshTokenRepository, times(2)).save(any(RefreshTokenEntity.class));
    }

    @Test
    @DisplayName("Test revoke by user with no tokens")
    void givenUserWithNoTokens_whenRevokeByUser_thenNothingHappens() {
        // given
        BDDMockito.given(refreshTokenRepository.findAllByUser(testUser)).willReturn(Collections.emptyList());

        // when
        refreshTokenService.revokeByUser(testUser);

        // then
        verify(refreshTokenRepository, never()).save(any());
    }
}
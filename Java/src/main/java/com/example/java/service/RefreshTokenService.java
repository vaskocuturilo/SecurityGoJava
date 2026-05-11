package com.example.java.service;

import com.example.java.dto.RefreshTokenRequest;
import com.example.java.entity.RefreshTokenEntity;
import com.example.java.entity.UserEntity;
import com.example.java.exception.UserException;
import com.example.java.repository.RefreshTokenRepository;
import com.example.java.repository.UserRepository;
import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class RefreshTokenService implements IRefreshTokenService {

    @Value("${security.jwt.refresh-token.expiration:604800000}")
    private long refreshTokenExpiration;

    private final RefreshTokenRepository refreshTokenRepository;

    private final UserRepository userRepository;

    @Override
    @Transactional
    public RefreshTokenEntity createRefreshToken(String email) {
        final UserEntity existUser = userRepository.findByEmail(email)
                .orElseThrow(() -> new UserException("Unknow user", HttpStatus.NOT_FOUND));

        refreshTokenRepository.deleteByUser(existUser);

        final RefreshTokenEntity refreshToken = RefreshTokenEntity
                .builder()
                .token(UUID.randomUUID().toString())
                .user(existUser)
                .expiresAt(LocalDateTime.now().plus(refreshTokenExpiration, ChronoUnit.MILLIS))
                .revoked(false)
                .build();

        return refreshTokenRepository.save(refreshToken);
    }

    @Override
    public RefreshTokenEntity verifyExpiration(String rawToken) {
        final RefreshTokenEntity token = refreshTokenRepository
                .findByToken(rawToken)
                .orElseThrow(() -> new UserException("Invalid refresh token",
                        HttpStatus.UNAUTHORIZED));

        if (token.isRevoked()) {
            throw new UserException("Refresh token was revoked", HttpStatus.UNAUTHORIZED);
        }

        if (token.getExpiresAt().isBefore(LocalDateTime.now())) {
            refreshTokenRepository.delete(token);
            throw new UserException("Refresh token has expired, please log in again",
                    HttpStatus.UNAUTHORIZED);
        }

        return token;
    }

    @Transactional
    public void logout(final RefreshTokenRequest request) {
        final RefreshTokenEntity token = refreshTokenRepository
                .findByToken(request.refreshToken())
                .orElseThrow(() -> new UserException("Invalid refresh token",
                        HttpStatus.UNAUTHORIZED));

        token.setRevoked(true);
        refreshTokenRepository.save(token);
    }

    @Transactional
    public void revokeByUser(UserEntity user) {
        refreshTokenRepository.findAllByUser(user)
                .forEach(token -> {
                    token.setRevoked(true);
                    refreshTokenRepository.save(token);
                });
    }
}

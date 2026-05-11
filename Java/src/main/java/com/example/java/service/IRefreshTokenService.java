package com.example.java.service;

import com.example.java.dto.RefreshTokenRequest;
import com.example.java.entity.RefreshTokenEntity;
import com.example.java.entity.UserEntity;

public interface IRefreshTokenService {

    RefreshTokenEntity createRefreshToken(String email);

    RefreshTokenEntity verifyExpiration(String rawToken);

    void revokeByUser(UserEntity user);

    void logout(final RefreshTokenRequest request);
}

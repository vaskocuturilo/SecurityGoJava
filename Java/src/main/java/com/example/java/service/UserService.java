package com.example.java.service;

import com.example.java.config.UserAuthenticationProvider;
import com.example.java.dto.*;
import com.example.java.entity.RefreshTokenEntity;
import com.example.java.entity.UserEntity;
import com.example.java.exception.UserException;
import com.example.java.mapper.UserMapper;
import com.example.java.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.nio.CharBuffer;
import java.util.Arrays;

@Service
@RequiredArgsConstructor
public class UserService implements IUserService {

    private final UserMapper userMapper;

    private final UserRepository userRepository;

    private final PasswordEncoder passwordEncoder;

    private final RefreshTokenService refreshTokenService;

    private final UserAuthenticationProvider userAuthenticationProvider;

    private static final String UNKNOWN_USER = "Unknown user";


    public AuthResponse login(final CredentialsDto credentials) {
        UserEntity user = userRepository.findByEmail(credentials.email())
                .orElseThrow(() -> new UserException(UNKNOWN_USER, HttpStatus.NOT_FOUND));

        if (!user.isEnabled()) {
            throw new UserException("The user is not enabled", HttpStatus.FORBIDDEN);
        }

        try {
            if (!passwordEncoder.matches(CharBuffer.wrap(credentials.password()), user.getPassword())) {
                throw new UserException("Invalid password", HttpStatus.BAD_REQUEST);
            }
            final UserDto userDto = userMapper.toUserDto(user);
            final String accessToken = userAuthenticationProvider.createToken(userDto);
            final String refreshToken = refreshTokenService.createRefreshToken(user.getEmail()).getToken();

            return new AuthResponse(userDto, accessToken, refreshToken);

        } finally {
            Arrays.fill(credentials.password(), '\0');
        }
    }

    public AuthResponse register(final SignUpDto userDto) {
        if (userRepository.findByEmail(userDto.email()).isPresent()) {
            throw new UserException("Email already exists", HttpStatus.BAD_REQUEST);
        }

        final UserEntity user = new UserEntity();

        try {
            user.setPassword(passwordEncoder.encode(CharBuffer.wrap(userDto.password())));
        } finally {
            Arrays.fill(userDto.password(), '\0');
        }

        user.setEmail(userDto.email());

        user.setUsername(userDto.email());
        user.setEnabled(true);

        final UserEntity savedUser = userRepository.save(user);
        final UserDto savedUserDto = userMapper.toUserDto(savedUser);
        final String accessToken = userAuthenticationProvider.createToken(savedUserDto);
        final String refreshToken = refreshTokenService.createRefreshToken(savedUser.getEmail()).getToken();

        return new AuthResponse(savedUserDto, accessToken, refreshToken);
    }

    public AuthResponse refresh(final RefreshTokenRequest request) {
        final RefreshTokenEntity refreshToken = refreshTokenService.verifyExpiration(request.refreshToken());

        final UserDto userDto = userMapper.toUserDto(refreshToken.getUser());
        final String newAccessToken = userAuthenticationProvider.createToken(userDto);
        final String newRefreshToken = refreshTokenService
                .createRefreshToken(userDto.email())
                .getToken();

        return new AuthResponse(userDto, newAccessToken, newRefreshToken);
    }

    public void logoutAll(final String email) {
        final UserEntity user = userRepository.findByEmail(email).orElseThrow(() -> new UserException("Unknown user", HttpStatus.NOT_FOUND));
        refreshTokenService.revokeByUser(user);
    }

    public UserDto findByLogin(final String email) {
        final UserEntity user = userRepository.findByEmail(email)
                .orElseThrow(() -> new UserException(UNKNOWN_USER, HttpStatus.NOT_FOUND));
        return userMapper.toUserDto(user);
    }
}

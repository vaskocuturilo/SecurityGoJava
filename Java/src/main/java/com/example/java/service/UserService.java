package com.example.java.service;

import com.example.java.dto.CredentialsDto;
import com.example.java.dto.SignUpDto;
import com.example.java.dto.UserDto;
import com.example.java.entity.UserEntity;
import com.example.java.exception.UserException;
import com.example.java.mapper.UserMapper;
import com.example.java.repository.UserRepository;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.nio.CharBuffer;
import java.time.LocalDateTime;

@Service
public class UserService implements IUserService {

    private final UserMapper userMapper;

    private final UserRepository userRepository;

    private final PasswordEncoder passwordEncoder;

    private static final String UNKNOWN_USER = "Unknown user";

    public UserService(UserMapper userMapper,
                       UserRepository userRepository,
                       PasswordEncoder passwordEncoder) {
        this.userMapper = userMapper;
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
    }

    public UserDto login(final CredentialsDto credentials) {
        UserEntity user = userRepository.findByEmail(credentials.email())
                .orElseThrow(() -> new UserException(UNKNOWN_USER, HttpStatus.NOT_FOUND));

        if (!user.isEnabled()) {
            throw new UserException("The user is not enabled", HttpStatus.FORBIDDEN);
        }

        if (passwordEncoder.matches(CharBuffer.wrap(credentials.password()), user.getPassword())) {
            return userMapper.toUserDto(user);
        }

        throw new UserException("Invalid password", HttpStatus.BAD_REQUEST);
    }

    public UserDto register(final SignUpDto userDto) {
        if (userRepository.findByEmail(userDto.email()).isPresent()) {
            throw new UserException("Email already exists", HttpStatus.BAD_REQUEST);
        }

        final UserEntity user = new UserEntity();

        user.setPassword(passwordEncoder.encode(CharBuffer.wrap(userDto.password())));

        user.setUsername(userDto.email());

        user.setEmail(userDto.email());

        user.setEnabled(true);

        user.setCreatedAt(LocalDateTime.now());

        user.setUpdatedAt(LocalDateTime.now());

        final UserEntity savedUser = userRepository.save(user);

        return userMapper.toUserDto(savedUser);
    }

    public UserDto findByLogin(final String login) {
        final UserEntity user = userRepository.findByEmail(login)
                .orElseThrow(() -> new UserException(UNKNOWN_USER, HttpStatus.NOT_FOUND));
        return userMapper.toUserDto(user);
    }
}

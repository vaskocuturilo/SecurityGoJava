package com.example.java.rest;


import com.example.java.config.UserAuthenticationProvider;
import com.example.java.dto.AuthResponse;
import com.example.java.dto.CredentialsDto;
import com.example.java.dto.SignUpDto;
import com.example.java.dto.UserDto;
import com.example.java.service.IUserService;
import jakarta.validation.Valid;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.net.URI;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/users")
public class UserRestControllerV1 {

    private final IUserService userService;
    private final UserAuthenticationProvider userAuthenticationProvider;

    @PostMapping("/login")
    public ResponseEntity<@NonNull AuthResponse> login(@RequestBody @Valid CredentialsDto credentials) {
        final UserDto userDto = userService.login(credentials);
        final String token = userAuthenticationProvider.createToken(userDto);

        return ResponseEntity.ok(new AuthResponse(userDto, token));
    }

    @PostMapping("/register")
    public ResponseEntity<@NonNull AuthResponse> register(@RequestBody @Valid SignUpDto userSignUp) {
        final UserDto createdUser = userService.register(userSignUp);
        final String token = userAuthenticationProvider.createToken(createdUser);
        return ResponseEntity
                .created(URI.create("/users/" + createdUser.id()))
                .body(new AuthResponse(createdUser, token));
    }

    @PostMapping("/logout")
    public ResponseEntity<@NonNull Void> logout() {
        return ResponseEntity.ok().build();
    }
}

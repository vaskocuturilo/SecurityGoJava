package com.example.java.config;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.example.java.dto.UserDto;
import jakarta.annotation.PostConstruct;
import lombok.Getter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Component;

import java.util.Base64;
import java.util.Collections;
import java.util.Date;

@Getter
@Component
public class UserAuthenticationProvider {

    @Value("${security.jwt.token.secret-key:secret-key}")
    private String secretKey;

    @Value("${security.jwt.token.expiration:300000}")
    private long expiration;

    @PostConstruct
    protected void init() {
        secretKey = Base64.getEncoder().encodeToString(secretKey.getBytes());
    }

    public String createToken(UserDto user) {
        final Date now = new Date();
        final Date validity = new Date(now.getTime() + expiration);

        final Algorithm algorithm = Algorithm.HMAC256(secretKey);
        return JWT.create()
                .withSubject(user.email())
                .withIssuedAt(now)
                .withExpiresAt(validity)
                .withClaim("username", user.username())
                .sign(algorithm);
    }

    public Authentication validateToken(String token) {
        final Algorithm algorithm = Algorithm.HMAC256(secretKey);
        final DecodedJWT decoded = JWT.require(algorithm).build().verify(token);

        final UserDto user = UserDto.builder()
                .email(decoded.getSubject())
                .username(decoded.getClaim("username").asString())
                .build();

        return new UsernamePasswordAuthenticationToken(user, null, Collections.emptyList());
    }

    public String extractEmail(String token) {
        final Algorithm algorithm = Algorithm.HMAC256(secretKey);
        return JWT.require(algorithm)
                .build()
                .verify(token)
                .getSubject();
    }
}
package com.example.java.rest;

import com.example.java.AbstractRestControllerBaseTest;
import com.example.java.dto.CredentialsDto;
import com.example.java.dto.RefreshTokenRequest;
import com.example.java.dto.SignUpDto;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import tools.jackson.databind.ObjectMapper;


import static org.hamcrest.Matchers.containsString;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("test")
class UserRestControllerV1Test extends AbstractRestControllerBaseTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @Value("${http.auth-token-header-name}")
    private String headerName;

    @Value("${http.auth-token}")
    private String authToken;

    private static final String BASE_URL = "/api/v1/users";
    private static final String TEST_EMAIL = "test@test.com";
    private static final String TEST_PASSWORD = "password123";
    private static final String TEST_REFRESH_TOKEN = "mocked-refresh-token";

    @Test
    @DisplayName("login: valid credentials return 200 with AuthResponse")
    void givenValidCredentials_whenLogin_thenReturn200() throws Exception {
        final CredentialsDto credentials = new CredentialsDto(
                TEST_EMAIL, TEST_PASSWORD.toCharArray());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(credentials)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.accessToken").isNotEmpty())
                .andExpect(jsonPath("$.refreshToken").isNotEmpty())
                .andExpect(jsonPath("$.user.email").value(TEST_EMAIL));
    }

    @Test
    @DisplayName("login: unknown email returns 404")
    void givenUnknownEmail_whenLogin_thenReturn404() throws Exception {
        final CredentialsDto credentials = new CredentialsDto(
                "unknown@test.com", TEST_PASSWORD.toCharArray());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(credentials)))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.error").isNotEmpty());
    }

    @Test
    @DisplayName("login: wrong password returns 400")
    void givenWrongPassword_whenLogin_thenReturn400() throws Exception {
        final CredentialsDto credentials = new CredentialsDto(
                TEST_EMAIL, "wrongpassword".toCharArray());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(credentials)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error").isNotEmpty());
    }

    @Test
    @DisplayName("login: invalid email format returns 400")
    void givenInvalidEmail_whenLogin_thenReturn400() throws Exception {
        final CredentialsDto credentials = new CredentialsDto(
                "not-an-email", TEST_PASSWORD.toCharArray());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(credentials)))
                .andExpect(status().isBadRequest());
    }

    @Test
    @DisplayName("login: missing API key returns 401")
    void givenMissingApiKey_whenLogin_thenReturn401() throws Exception {
        final CredentialsDto credentials = new CredentialsDto(
                TEST_EMAIL, TEST_PASSWORD.toCharArray());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(credentials)))
                .andExpect(status().isUnauthorized());
    }

    @Test
    @DisplayName("register: new user returns 201 with AuthResponse")
    void givenNewUser_whenRegister_thenReturn201() throws Exception {
        final SignUpDto signUpDto = new SignUpDto(
                "newuser@test.com", TEST_PASSWORD.toCharArray());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(signUpDto)))
                .andExpect(status().isCreated())
                .andExpect(header().string("Location", containsString("/users/")))
                .andExpect(jsonPath("$.accessToken").isNotEmpty())
                .andExpect(jsonPath("$.refreshToken").isNotEmpty())
                .andExpect(jsonPath("$.user.email").value("newuser@test.com"));
    }

    @Test
    @DisplayName("register: duplicate email returns 400")
    void givenDuplicateEmail_whenRegister_thenReturn400() throws Exception {
        final SignUpDto signUpDto = new SignUpDto(
                TEST_EMAIL, TEST_PASSWORD.toCharArray());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(signUpDto)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error").value("Email already exists"));
    }

    @Test
    @DisplayName("register: invalid email format returns 400")
    void givenInvalidEmail_whenRegister_thenReturn400() throws Exception {
        final SignUpDto signUpDto = new SignUpDto(
                "not-an-email", TEST_PASSWORD.toCharArray());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(signUpDto)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.errors.email").value("Email is not valid"));
    }

    @Test
    @DisplayName("register: missing API key returns 401")
    void givenMissingApiKey_whenRegister_thenReturn401() throws Exception {
        final SignUpDto signUpDto = new SignUpDto(
                "newuser@test.com", TEST_PASSWORD.toCharArray());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(signUpDto)))
                .andExpect(status().isUnauthorized());
    }

    @Test
    @DisplayName("refresh: valid refresh token returns 200 with new AuthResponse")
    void givenValidRefreshToken_whenRefresh_thenReturn200() throws Exception {
        final CredentialsDto credentials = new CredentialsDto(
                TEST_EMAIL, TEST_PASSWORD.toCharArray());

        final String loginResponse = mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(credentials)))
                .andExpect(status().isOk())
                .andReturn()
                .getResponse()
                .getContentAsString();

        final String refreshToken = objectMapper.readTree(loginResponse)
                .get("refreshToken")
                .asText();

        final RefreshTokenRequest request = new RefreshTokenRequest(refreshToken);

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/refresh")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.accessToken").isNotEmpty())
                .andExpect(jsonPath("$.refreshToken").isNotEmpty())
                .andExpect(jsonPath("$.user.email").value(TEST_EMAIL));
    }

    @Test
    @DisplayName("refresh: invalid refresh token returns 401")
    void givenInvalidRefreshToken_whenRefresh_thenReturn401() throws Exception {
        final RefreshTokenRequest request = new RefreshTokenRequest("invalid-token");

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/refresh")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isUnauthorized())
                .andExpect(jsonPath("$.error").value("Invalid refresh token"));
    }

    @Test
    @DisplayName("refresh: missing API key returns 401")
    void givenMissingApiKey_whenRefresh_thenReturn401() throws Exception {
        final RefreshTokenRequest request = new RefreshTokenRequest(TEST_REFRESH_TOKEN);

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/refresh")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isUnauthorized());
    }

    @Test
    @DisplayName("logout: valid refresh token returns 200")
    void givenValidRefreshToken_whenLogout_thenReturn200() throws Exception {
        final CredentialsDto credentials = new CredentialsDto(
                TEST_EMAIL, TEST_PASSWORD.toCharArray());

        final String loginResponse = mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(credentials)))
                .andExpect(status().isOk())
                .andReturn()
                .getResponse()
                .getContentAsString();

        final String refreshToken = objectMapper.readTree(loginResponse)
                .get("refreshToken")
                .asText();

        final RefreshTokenRequest request = new RefreshTokenRequest(refreshToken);

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/logout")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk());
    }

    @Test
    @DisplayName("logout: invalid refresh token returns 401")
    void givenInvalidRefreshToken_whenLogout_thenReturn401() throws Exception {
        final RefreshTokenRequest request = new RefreshTokenRequest("invalid-token");

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/logout")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isUnauthorized())
                .andExpect(jsonPath("$.error").value("Invalid refresh token"));
    }

    @Test
    @DisplayName("logout: token is unusable after logout")
    void givenValidRefreshToken_whenLogoutThenRefresh_thenReturn401() throws Exception {
        final CredentialsDto credentials = new CredentialsDto(
                TEST_EMAIL, TEST_PASSWORD.toCharArray());

        final String loginResponse = mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(credentials)))
                .andReturn().getResponse().getContentAsString();

        final String refreshToken = objectMapper.readTree(loginResponse)
                .get("refreshToken").asText();

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/logout")
                        .contentType(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(
                                new RefreshTokenRequest(refreshToken))))
                .andExpect(status().isOk());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/refresh")
                        .contentType(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(
                                new RefreshTokenRequest(refreshToken))))
                .andExpect(status().isUnauthorized())
                .andExpect(jsonPath("$.error").value("Refresh token was revoked"));
    }

    @Test
    @DisplayName("logoutAll: authenticated user returns 200 and all tokens revoked")
    void givenAuthenticatedUser_whenLogoutAll_thenReturn200() throws Exception {
        final CredentialsDto credentials = new CredentialsDto(
                TEST_EMAIL, TEST_PASSWORD.toCharArray());

        final String loginResponse = mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(credentials)))
                .andExpect(status().isOk())
                .andReturn().getResponse().getContentAsString();

        final String accessToken = objectMapper.readTree(loginResponse)
                .get("accessToken").asText();
        final String refreshToken = objectMapper.readTree(loginResponse)
                .get("refreshToken").asText();

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/logout-all")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .header(HttpHeaders.AUTHORIZATION, "Bearer " + accessToken))
                .andExpect(status().isOk());

        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/refresh")
                        .contentType(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken)
                        .content(objectMapper.writeValueAsString(
                                new RefreshTokenRequest(refreshToken))))
                .andExpect(status().isUnauthorized());
    }

    @Test
    @DisplayName("logoutAll: missing JWT returns 401")
    void givenMissingJwt_whenLogoutAll_thenReturn401() throws Exception {
        mockMvc.perform(MockMvcRequestBuilders
                        .post(BASE_URL + "/logout-all")
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON)
                        .header(headerName, authToken))
                .andExpect(status().isUnauthorized());
    }
}
package com.example.java.rest;

import com.example.java.AbstractRestControllerBaseTest;
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

import static org.hamcrest.Matchers.hasSize;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("test")
class RestTaskControllerV1Test extends AbstractRestControllerBaseTest {

    @Autowired
    private MockMvc mockMvc;

    @Value("${http.auth-token-header-name}")
    private String headerName;

    @Value("${http.auth-token}")
    private String authToken;

    private static final String ENDPOINT_PATH = "/api/v1/tasks";

    @Test
    void givenTasks_whenGetTasks_thenStatus200() throws Exception {
        mockMvc.perform(MockMvcRequestBuilders
                        .get(ENDPOINT_PATH)
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON).header(headerName, authToken)
                .header(HttpHeaders.AUTHORIZATION, "Bearer " + generateTestToken()))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$[*]").isNotEmpty())
                .andExpect(jsonPath("$[*]", hasSize(3)));
    }

    @Test
    void givenTasks_whenGetTasks_thenStatus401() throws Exception {
        mockMvc.perform(MockMvcRequestBuilders
                        .get(ENDPOINT_PATH)
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON).header(headerName, authToken))
                .andExpect(status().isUnauthorized());
    }

    @Test
    void givenWithoutAPikey_whenGetTasksWithoutApiKey_thenStatus403() throws Exception {
        mockMvc.perform(MockMvcRequestBuilders
                        .get(ENDPOINT_PATH)
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON))
                .andExpect(status().isForbidden());
    }
}
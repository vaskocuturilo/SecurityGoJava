package com.example.java.rest;

import com.example.java.AbstractRestControllerBaseTest;
import com.example.java.dto.Task;
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

    @Autowired
    private ObjectMapper objectMapper;

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
                .andExpect(jsonPath("$[*]", hasSize(4)));
    }

    @Test
    void givenTasks_whenCreateTask_thenStatus200() throws Exception {
        final Task task = new Task("Task4", "Task4");

        mockMvc.perform(MockMvcRequestBuilders
                        .post(ENDPOINT_PATH)
                        .contentType(MediaType.APPLICATION_JSON)
                        .accept(MediaType.APPLICATION_JSON).header(headerName, authToken)
                        .header(HttpHeaders.AUTHORIZATION, "Bearer " + generateTestToken())
                        .content(objectMapper.writeValueAsString(task)))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$[*]").isNotEmpty())
                .andExpect(jsonPath("$[*]", hasSize(2)));
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
package com.example.java.service;

import com.example.java.dto.Task;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.junit.jupiter.api.Assertions.assertThrows;

@ExtendWith(MockitoExtension.class)
class TaskServiceTest {

    @InjectMocks
    private TaskService taskService;

    private static final String TEST_KEY = "task-key";
    private static final String TEST_VALUE = "task-value";



    @BeforeEach
    void setup() {
        taskService = new TaskService();
    }


    @Test
    @DisplayName("getAllTasks: returns empty map when no tasks exist")
    void givenNoTasks_whenGetAllTasks_thenReturnEmptyMap() {
        // when
        final Map<String, String> result = taskService.getTasks();

        // then
        assertThat(result).isNotNull();
    }

    @Test
    @DisplayName("getAllTasks: returns all existing tasks")
    void givenExistingTasks_whenGetAllTasks_thenReturnAllTasks() {
        // given
        taskService.createTask(new Task(TEST_KEY, TEST_VALUE));
        taskService.createTask(new Task("second-key", "second-value"));

        // when
        final Map<String, String> result = taskService.getTasks();

        // then
        assertThat(result).hasSize(5);
        assertThat(result).containsEntry(TEST_KEY, TEST_VALUE);
        assertThat(result).containsEntry("second-key", "second-value");
    }

    @Test
    @DisplayName("getAllTasks: returned map is unmodifiable")
    void givenExistingTasks_whenGetAllTasks_thenReturnUnmodifiableMap() {
        // given
        taskService.createTask(new Task(TEST_KEY, TEST_VALUE));

        // when
        final Map<String, String> result = taskService.getTasks();

        // then — attempt to modify must throw
        assertThatThrownBy(() -> result.put("new-key", "new-value"))
                .isInstanceOf(UnsupportedOperationException.class);
    }

    @Test
    @DisplayName("createTask: valid task is created and returned")
    void givenValidTask_whenCreateTask_thenTaskIsCreatedAndReturned() {
        // given
        final Task task = new Task(TEST_KEY, TEST_VALUE);

        // when
        final Task result = taskService.createTask(task);

        // then
        assertThat(result).isNotNull();
        assertThat(result.key()).isEqualTo(TEST_KEY);
        assertThat(result.value()).isEqualTo(TEST_VALUE);
    }

    @Test
    @DisplayName("createTask: created task is persisted and retrievable")
    void givenValidTask_whenCreateTask_thenTaskIsPersisted() {
        // given
        final Task task = new Task(TEST_KEY, TEST_VALUE);

        // when
        taskService.createTask(task);

        // then
        final Map<String, String> allTasks = taskService.getTasks();
        assertThat(allTasks).containsEntry(TEST_KEY, TEST_VALUE);
    }

    @Test
    @DisplayName("createTask: duplicate key throws UserException with 409")
    void givenDuplicateKey_whenCreateTask_thenThrowUserExceptionConflict() {
        // given
        taskService.createTask(new Task(TEST_KEY, TEST_VALUE));

        // when / then
        final IllegalArgumentException exception = assertThrows(IllegalArgumentException.class,
                () -> taskService.createTask(new Task(TEST_KEY, "different-value")));

        assertThat(exception.getMessage()).isEqualTo("Task with key already exists");
    }

    @Test
    @DisplayName("createTask: duplicate key does not overwrite existing value")
    void givenDuplicateKey_whenCreateTask_thenOriginalValueIsPreserved() {
        // given
        taskService.createTask(new Task(TEST_KEY, TEST_VALUE));

        // when
        assertThrows(IllegalArgumentException.class,
                () -> taskService.createTask(new Task(TEST_KEY, "different-value")));

        // then
        assertThat(taskService.getTasks()).containsEntry(TEST_KEY, TEST_VALUE);
    }

    @Test
    @DisplayName("createTask: multiple different tasks can be created")
    void givenMultipleDifferentTasks_whenCreateTask_thenAllAreCreated() {
        // given
        final Task first = new Task("key-1", "value-1");
        final Task second = new Task("key-2", "value-2");
        final Task third = new Task("key-3", "value-3");

        // when
        taskService.createTask(first);
        taskService.createTask(second);
        taskService.createTask(third);

        // then
        final Map<String, String> allTasks = taskService.getTasks();
        assertThat(allTasks).hasSize(6);
        assertThat(allTasks).containsEntry("key-1", "value-1");
        assertThat(allTasks).containsEntry("key-2", "value-2");
        assertThat(allTasks).containsEntry("key-3", "value-3");
    }

    @Test
    @DisplayName("createTask: tasks are isolated between test runs via setUp")
    void givenFreshService_whenCreateTask_thenOnlyOneTaskExists() {
        // given — setUp() creates a fresh TaskService instance for each test
        final Task task = new Task(TEST_KEY, TEST_VALUE);

        // when
        taskService.createTask(task);

        assertThat(taskService.getTasks()).hasSize(4);
    }
}
#!/usr/bin/env python3
# =============================================================================
# test_endpoints.py — Скрипт для тестирования API эндпоинтов
# =============================================================================
# Автоматическая проверка всех ключевых точек API
# =============================================================================

import requests
import json
import sys
import time
from typing import Dict, List, Tuple
from dataclasses import dataclass
from enum import Enum


# =============================================================================
# Константы
# =============================================================================

class Colors:
    """ANSI цвета для вывода."""
    GREEN = '\033[92m'
    RED = '\033[91m'
    YELLOW = '\033[93m'
    BLUE = '\033[94m'
    CYAN = '\033[96m'
    RESET = '\033[0m'
    BOLD = '\033[1m'


class Status(Enum):
    """Статусы тестов."""
    PASS = "PASS"
    FAIL = "FAIL"
    SKIP = "SKIP"


@dataclass
class TestResult:
    """Результат теста."""
    endpoint: str
    method: str
    status_code: int
    expected_code: int
    elapsed_ms: float
    status: Status
    message: str = ""


# =============================================================================
# Конфигурация
# =============================================================================

# Базовый URL API (можно переопределить через переменную окружения)
BASE_URL = "http://127.0.0.1:8080"

# Таймаут запросов (секунды)
TIMEOUT = 5

# Заголовки по умолчанию
HEADERS = {
    "Content-Type": "application/json",
    "Accept": "application/json"
}


# =============================================================================
# Тестовые сценарии
# =============================================================================

def get_test_scenarios() -> List[Dict]:
    """Возвращает список тестовых сценариев."""
    return [
        {
            "name": "Health Check",
            "method": "GET",
            "endpoint": "/health",
            "expected_code": 200,
            "description": "Проверка работоспособности сервиса"
        },
        {
            "name": "Get Users (Empty)",
            "method": "GET",
            "endpoint": "/users",
            "expected_code": 200,
            "description": "Получение пустого списка пользователей"
        },
        {
            "name": "Get Users with Pagination",
            "method": "GET",
            "endpoint": "/users?page=1&per_page=10",
            "expected_code": 200,
            "description": "Получение списка с пагинацией"
        },
        {
            "name": "Create User (Alice)",
            "method": "POST",
            "endpoint": "/add-user",
            "expected_code": 201,
            "body": {"name": "Alice"},
            "description": "Создание первого пользователя"
        },
        {
            "name": "Create User (Bob)",
            "method": "POST",
            "endpoint": "/add-user",
            "expected_code": 201,
            "body": {"name": "Bob"},
            "description": "Создание второго пользователя"
        },
        {
            "name": "Create User (Duplicate)",
            "method": "POST",
            "endpoint": "/add-user",
            "expected_code": 409,
            "body": {"name": "Alice"},
            "description": "Попытка создания дубликата"
        },
        {
            "name": "Create User (Empty Name)",
            "method": "POST",
            "endpoint": "/add-user",
            "expected_code": 400,
            "body": {"name": ""},
            "description": "Валидация пустого имени"
        },
        {
            "name": "Create User (Too Long)",
            "method": "POST",
            "endpoint": "/add-user",
            "expected_code": 400,
            "body": {"name": "a" * 101},
            "description": "Валидация длинного имени"
        },
        {
            "name": "Get User by ID (Exists)",
            "method": "GET",
            "endpoint": "/user/1",
            "expected_code": 200,
            "description": "Получение существующего пользователя"
        },
        {
            "name": "Get User by ID (Not Found)",
            "method": "GET",
            "endpoint": "/user/999",
            "expected_code": 404,
            "description": "Получение несуществующего пользователя"
        },
        {
            "name": "Update User",
            "method": "PUT",
            "endpoint": "/user/1",
            "expected_code": 200,
            "body": {"name": "Alice Updated"},
            "description": "Обновление пользователя"
        },
        {
            "name": "Update User (Not Found)",
            "method": "PUT",
            "endpoint": "/user/999",
            "expected_code": 404,
            "body": {"name": "Unknown"},
            "description": "Обновление несуществующего"
        },
        {
            "name": "Get Stats",
            "method": "GET",
            "endpoint": "/stats/active",
            "expected_code": 200,
            "description": "Получение статистики"
        },
        {
            "name": "Delete User",
            "method": "DELETE",
            "endpoint": "/user/2",
            "expected_code": 200,
            "description": "Удаление пользователя"
        },
        {
            "name": "Delete User (Not Found)",
            "method": "DELETE",
            "endpoint": "/user/999",
            "expected_code": 404,
            "description": "Удаление несуществующего"
        },
        {
            "name": "Invalid JSON",
            "method": "POST",
            "endpoint": "/add-user",
            "expected_code": 400,
            "body": "{invalid json}",
            "description": "Некорректный JSON",
            "raw_body": True
        },
        {
            "name": "Invalid Content-Type",
            "method": "POST",
            "endpoint": "/add-user",
            "expected_code": 415,
            "body": {"name": "Test"},
            "headers": {"Content-Type": "text/plain"},
            "description": "Неверный Content-Type"
        },
        {
            "name": "Invalid Method",
            "method": "PATCH",
            "endpoint": "/user/1",
            "expected_code": 405,
            "description": "Неподдерживаемый метод"
        }
    ]


# =============================================================================
# Функции тестирования
# =============================================================================

def run_test(scenario: Dict) -> TestResult:
    """Выполняет один тест."""
    method = scenario["method"]
    endpoint = scenario["endpoint"]
    expected_code = scenario["expected_code"]
    url = f"{BASE_URL}{endpoint}"
    
    # Подготовка запроса
    kwargs = {
        "timeout": TIMEOUT,
        "headers": HEADERS.copy()
    }
    
    # Добавляем кастомные заголовки если есть
    if "headers" in scenario:
        kwargs["headers"].update(scenario["headers"])
    
    # Добавляем тело запроса если есть
    if "body" in scenario:
        if scenario.get("raw_body", False):
            kwargs["data"] = scenario["body"]
        else:
            kwargs["json"] = scenario["body"]
    
    # Выполнение запроса
    start_time = time.time()
    try:
        response = requests.request(method, url, **kwargs)
        elapsed_ms = (time.time() - start_time) * 1000
        
        # Проверка статуса
        if response.status_code == expected_code:
            status = Status.PASS
            message = f"✓ Ожидаемый статус {expected_code}"
        else:
            status = Status.FAIL
            message = f"✗ Ожидался {expected_code}, получен {response.status_code}"
        
        return TestResult(
            endpoint=endpoint,
            method=method,
            status_code=response.status_code,
            expected_code=expected_code,
            elapsed_ms=elapsed_ms,
            status=status,
            message=message
        )
        
    except requests.exceptions.ConnectionError as e:
        return TestResult(
            endpoint=endpoint,
            method=method,
            status_code=0,
            expected_code=expected_code,
            elapsed_ms=0,
            status=Status.FAIL,
            message=f"✗ Ошибка соединения: {BASE_URL} недоступен"
        )
    except requests.exceptions.Timeout:
        return TestResult(
            endpoint=endpoint,
            method=method,
            status_code=0,
            expected_code=expected_code,
            elapsed_ms=0,
            status=Status.FAIL,
            message=f"✗ Таймаут запроса ({TIMEOUT}s)"
        )
    except Exception as e:
        return TestResult(
            endpoint=endpoint,
            method=method,
            status_code=0,
            expected_code=expected_code,
            elapsed_ms=0,
            status=Status.FAIL,
            message=f"✗ Ошибка: {str(e)}"
        )


def print_header(text: str):
    """Печатает заголовок."""
    print(f"\n{Colors.CYAN}{Colors.BOLD}{'=' * 70}{Colors.RESET}")
    print(f"{Colors.CYAN}{Colors.BOLD}{text.center(70)}{Colors.RESET}")
    print(f"{Colors.CYAN}{Colors.BOLD}{'=' * 70}{Colors.RESET}\n")


def print_result(result: TestResult, scenario_name: str):
    """Печатает результат теста."""
    status_color = {
        Status.PASS: Colors.GREEN,
        Status.FAIL: Colors.RED,
        Status.SKIP: Colors.YELLOW
    }
    
    status_str = f"[{result.status.value}]"
    color = status_color.get(result.status, Colors.RESET)
    
    print(f"{color}{status_str}{Colors.RESET} {Colors.BOLD}{scenario_name}{Colors.RESET}")
    print(f"       {Colors.BLUE}{result.method}{Colors.RESET} {result.endpoint}")
    print(f"       Статус: {result.status_code} (ожидался {result.expected_code})")
    print(f"       Время: {result.elapsed_ms:.2f}ms")
    if result.message:
        print(f"       {result.message}")
    print()


def print_summary(results: List[TestResult]):
    """Печатает сводку."""
    total = len(results)
    passed = sum(1 for r in results if r.status == Status.PASS)
    failed = sum(1 for r in results if r.status == Status.FAIL)
    skipped = sum(1 for r in results if r.status == Status.SKIP)
    
    avg_time = sum(r.elapsed_ms for r in results if r.elapsed_ms > 0) / max(passed + failed, 1)
    
    print_header("СВОДКА")
    
    print(f"Всего тестов:    {total}")
    print(f"{Colors.GREEN}Пройдено:        {passed}{Colors.RESET}")
    print(f"{Colors.RED}Провалено:       {failed}{Colors.RESET}")
    print(f"{Colors.YELLOW}Пропущено:       {skipped}{Colors.RESET}")
    print(f"Среднее время:   {avg_time:.2f}ms")
    print()
    
    # Процент успеха
    success_rate = (passed / total * 100) if total > 0 else 0
    print(f"Успешность:      {success_rate:.1f}%")
    
    # Визуальная шкала
    bar_length = 40
    filled = int(bar_length * passed / total) if total > 0 else 0
    bar = "█" * filled + "░" * (bar_length - filled)
    print(f"[{Colors.GREEN}{bar}{Colors.RESET}]")
    print()


# =============================================================================
# Основная функция
# =============================================================================

def main():
    """Основная функция."""
    print_header("🧪 ТЕСТИРОВАНИЕ API ENDPOINTS")
    
    print(f"Базовый URL: {Colors.BLUE}{BASE_URL}{Colors.RESET}")
    print(f"Таймаут:     {TIMEOUT}s")
    print()
    
    # Проверка доступности API
    print(f"{Colors.YELLOW}Проверка доступности API...{Colors.RESET}")
    try:
        response = requests.get(f"{BASE_URL}/health", timeout=3)
        if response.status_code == 200:
            print(f"{Colors.GREEN}✓ API доступен{Colors.RESET}\n")
        else:
            print(f"{Colors.RED}✗ API вернул статус {response.status_code}{Colors.RESET}\n")
    except requests.exceptions.ConnectionError:
        print(f"{Colors.RED}✗ API недоступен по адресу {BASE_URL}{Colors.RESET}")
        print(f"\n{Colors.YELLOW}Убедитесь что сервер запущен:{Colors.RESET}")
        print(f"  go run cmd/server/main.go")
        print(f"  или")
        print(f"  docker-compose up -d\n")
        sys.exit(1)
    
    # Получение сценариев
    scenarios = get_test_scenarios()
    print(f"Запланировано тестов: {len(scenarios)}\n")
    
    # Выполнение тестов
    results: List[TestResult] = []
    
    for i, scenario in enumerate(scenarios, 1):
        print(f"{Colors.BOLD}Тест {i}/{len(scenarios)}: {scenario['description']}{Colors.RESET}")
        result = run_test(scenario)
        print_result(result, scenario["name"])
        results.append(result)
        
        # Небольшая задержка между запросами
        time.sleep(0.1)
    
    # Вывод сводки
    print_summary(results)
    
    # Вывод детальной статистики по методам
    print_header("СТАТИСТИКА ПО МЕТОДАМ")
    
    methods = {}
    for r in results:
        if r.method not in methods:
            methods[r.method] = {"pass": 0, "fail": 0}
        if r.status == Status.PASS:
            methods[r.method]["pass"] += 1
        else:
            methods[r.method]["fail"] += 1
    
    for method, stats in sorted(methods.items()):
        total = stats["pass"] + stats["fail"]
        print(f"{Colors.BLUE}{method}{Colors.RESET}: {stats['pass']}/{total} passed")
    
    print()
    
    # Возврат кода выхода
    failed = sum(1 for r in results if r.status == Status.FAIL)
    if failed > 0:
        print(f"{Colors.RED}❌ Тесты провалены: {failed} ошибок{Colors.RESET}\n")
        sys.exit(1)
    else:
        print(f"{Colors.GREEN}✅ Все тесты пройдены!{Colors.RESET}\n")
        sys.exit(0)


# =============================================================================
# Запуск
# =============================================================================

if __name__ == "__main__":
    main()

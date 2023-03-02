class TableHeader:
    def __init__(self, names: list) -> None:
        self._names = names

    def __str__(self) -> str:
        return f"| {' | '.join(self._names)} |\n| {'------------ |' + '------------: |' * (len(self._names) - 1)}"


class ShortTableRaw:
    # rph -> requests per hour
    def __init__(self, name: str, average_rps: float, peak_rps: float) -> None:
        self._name = name
        self._average_rps = average_rps
        self._peak_rps = peak_rps

    def get_average_rps(self) -> float:
        return self._average_rps

    def get_peak_rps(self) -> float:
        return self._peak_rps

    def get_header(self) -> TableHeader:
        return TableHeader([
            "Requirement", "Average RPS", "Peak RPS",
        ])

    def __str__(self) -> str:
        return f"| {self._name} | {self._average_rps:.2f} | {self._peak_rps:.2f} |"


class TableRaw:
    # rph -> requests per hour
    def __init__(self, name: str, average_rph_from_user: float, peak_rph_from_user: float, users_count: int) -> None:
        self._name = name
        self._average_rph_from_user = average_rph_from_user
        self._peak_rph_from_user = peak_rph_from_user
        SECONDS_IN_HOUR = 3600
        self._average_rps = users_count * average_rph_from_user / SECONDS_IN_HOUR
        self._peak_rps = users_count * peak_rph_from_user / SECONDS_IN_HOUR

    def get_average_rps(self) -> float:
        return self._average_rps

    def get_peak_rps(self) -> float:
        return self._peak_rps

    def get_header(self) -> TableHeader:
        return TableHeader([
            "Requirement", "Average load (requests per hour from one user)",
            "Peak load (requests per hour from one user)", "Average RPS", "Peak RPS",
        ])

    def __str__(self) -> str:
        return f"| {self._name} | {self._average_rph_from_user:.2f} | " + \
            f"{self._peak_rph_from_user:.2f} | {self._average_rps:.2f} | {self._peak_rps:.2f} |"


class Message:
    def __init__(self, min_size: int, max_size: int, bytes_per_symbol: int) -> None:
        self._min = min_size * bytes_per_symbol
        self._max = max_size * bytes_per_symbol

    def get_min_bytes(self) -> int:
        return self._min

    def get_max_bytes(self) -> int:
        return self._max


class ComplesMessage:
    def __init__(self, messages: list) -> None:
        self._min = sum([el.get_min_bytes() for el in messages])
        self._max = sum([el.get_max_bytes() for el in messages])

    def get_min_bytes(self) -> int:
        return self._min

    def get_max_bytes(self) -> int:
        return self._max


class HandlerRaw:
    def __init__(self, name: str, requirements: list, request: Message, response: Message) -> None:
        self._name = name
        self._average_rps = sum([r.get_average_rps() for r in requirements])
        self._peak_rps = sum([r.get_peak_rps() for r in requirements])
        self._request = request
        self._response = response

    def get_average_rps(self) -> float:
        return self._average_rps

    def get_peak_rps(self) -> float:
        return self._peak_rps

    def get_request(self) -> Message:
        return self._request

    def get_response(self) -> Message:
        return self._request

    def get_header(self) -> TableHeader:
        return TableHeader([
            "Handler", "Average RPS", "Peak RPS", "Request MAX bytes", "Response MAX bytes",
        ])

    def __str__(self) -> str:
        return f"| {self._name} | {self._average_rps:.2f} | {self._peak_rps:.2f} | " + \
            f"{self._request.get_max_bytes()} | {self._response.get_max_bytes()} |"


class DatabaseRaw:
    def __init__(self, read_handlers, write_handlers, content_producer_handlers) -> None:
        self._read_average_rps = sum([el.get_average_rps() for el in read_handlers])
        self._read_peak_rps = sum([el.get_peak_rps() for el in read_handlers])
        self._write_average_rps = sum([el.get_average_rps() for el in write_handlers])
        self._write_peak_rps = sum([el.get_peak_rps() for el in write_handlers])

        self._average_rw_ratio = self._read_average_rps / self._write_average_rps * 100
        self._peak_rw_ratio = self._read_peak_rps / self._write_peak_rps * 100

        content_producer_average_rps = sum([el.get_average_rps() for el in content_producer_handlers])
        database_growth_per_request = sum([el.get_request().get_max_bytes() for el in content_producer_handlers])

        SECONDS_IN_DAY = 60 * 60 * 24
        SECONDS_IN_MONTH = SECONDS_IN_DAY * 30

        BYTES_IN_GIGABYTE = 1073741824
        self._day_database_growth_gb = database_growth_per_request * content_producer_average_rps * SECONDS_IN_DAY / BYTES_IN_GIGABYTE
        self._month_database_growth_gb = database_growth_per_request * content_producer_average_rps * SECONDS_IN_MONTH / BYTES_IN_GIGABYTE

    def get_header(self) -> TableHeader:
        return TableHeader([
            "Average R/W ratio", "Agerage read RPS", "Agerage write RPS",
            "Peak R/W ratio", "Peak read RPS", "Peak write RPS",
            "Database growth GB (per day)", "Database growth GB (per month)",
        ])


    def __str__(self) -> str:
        return f"| {self._average_rw_ratio:.2f}% | {self._read_average_rps:.2f} | {self._write_average_rps:.2f} | " + \
            f"{self._peak_rw_ratio:.2f}% | {self._read_peak_rps:.2f} | {self._write_peak_rps:.2f} | " + \
            f"{self._day_database_growth_gb:.2f} | {self._month_database_growth_gb:.2f} |"


def print_all(lst: list, header: str | None = None):
    if header:
        print(f"## {header}")
    if len(lst) > 0:
        print(lst[0].get_header())
    for el in lst:
        print(el)
    print()

PEAK_HOUR_RPS = 10_000
DAU_USERS = 50_000
MAU_USERS = 500_000
OVERALL_USERS = 1_000_000


# Auth actions
registration = ShortTableRaw("Registration", 2.0, 100.0)
authorization = ShortTableRaw("Authorization", 5.0, 200.0)
print_all([registration, authorization], "Auth")

# Actions with users
find_user_by_nickname = TableRaw("Find a user by nickname", 1.0, 10.0, DAU_USERS)
show_profile = TableRaw("Show profile", 0.5, 10.0, DAU_USERS)
print_all([find_user_by_nickname, show_profile], "Actions with users")

# Actions with own profile
update_profile_info = TableRaw("Update profile information", 0.1, 2.0, DAU_USERS)
set_profile_photo = TableRaw("Set a profile photo (select from preset)", 0.1, 2.0, DAU_USERS)
print_all([find_user_by_nickname, show_profile], "Actions with own profile")

# Actions with dialogs
create_dialog = TableRaw("Create dialog", 0.1, 10.0, DAU_USERS)
get_user_dialogs = TableRaw("Get user dialogs", 15.00, 100.0, DAU_USERS)
delete_dialog = TableRaw("Delete dialog", 0.01, 2.0, DAU_USERS)
send_message_to_dialog = TableRaw("Send a message to the dialog", 5.0, 100.0, DAU_USERS)
delete_message_from_dialog = TableRaw("Delete a message from a dialog (for both)", 0.1, 5.0, DAU_USERS)
find_message_in_dialog = TableRaw("Find a message in the dialog", 0.5, 5.0, DAU_USERS)
search_message_in_all_dialogs = TableRaw("Search for a message in all user dialogs", 1.0, 10.0, DAU_USERS)
view_all_sandbox_links_in_dialog = TableRaw("View all links to the sandbox in the dialogs", 0.2, 4.0, DAU_USERS)
print_all([create_dialog, delete_dialog, send_message_to_dialog, delete_message_from_dialog,
            find_message_in_dialog, search_message_in_all_dialogs, view_all_sandbox_links_in_dialog],
            "Actions with dialogs")


# Actions with sandbox
create_code_listing = TableRaw("Create code listing", 0.5, 1.5, DAU_USERS)
find_code_listing = TableRaw("Find code listing", 3.0, 15.0, DAU_USERS)
update_code_listing = TableRaw("Update code listing", 0.6, 3.0, DAU_USERS)
run_go_code = TableRaw("Run code (only Go)", 0.25, 2.0, DAU_USERS)
lint_go_code = TableRaw("Linting code (only Go)", 0.25, 1.5, DAU_USERS)
print_all([find_code_listing, create_code_listing, update_code_listing, run_go_code, lint_go_code],
          "Actions with sandbox",)

# common
BYTES_PER_SYMBOL = 2
void_message = Message(1, 2, BYTES_PER_SYMBOL)

# User service
user_data_message = Message(10 , 512, BYTES_PER_SYMBOL)
credential_message = Message(10 , 512, BYTES_PER_SYMBOL)
token_message = Message(150, 1000, BYTES_PER_SYMBOL)
user_id_message = Message(12, 20, BYTES_PER_SYMBOL)
create_user_request = ComplesMessage([user_data_message, credential_message])

create_user = HandlerRaw("CreateUser", [registration], create_user_request, user_id_message)
# TODO: add required handlers
get_user_by_id = HandlerRaw("GetUserByID", [authorization, find_user_by_nickname], user_id_message, user_data_message)
delete_user_by_id = HandlerRaw("DeleteUserByID", [], user_id_message, void_message)
generate_token = HandlerRaw("GenerateToken", [authorization], credential_message, token_message)
print_all([create_user, get_user_by_id, delete_user_by_id, generate_token], "User service")

database_stat = DatabaseRaw([get_user_by_id, generate_token], [create_user], [create_user])
print_all([database_stat], "User service databse load")


# Sandbox service
pogram_sourse_message = Message(1 , 4000, BYTES_PER_SYMBOL)
pogram_id_message = Message(12, 20, BYTES_PER_SYMBOL)
program_output_message = Message(1 , 4000, BYTES_PER_SYMBOL)
program_message = ComplesMessage([pogram_id_message, pogram_sourse_message, program_output_message, program_output_message])
update_program_sourse_request_message = ComplesMessage([pogram_id_message, pogram_sourse_message])

get_program_by_id = HandlerRaw("GetProgramByID", [find_code_listing], pogram_id_message, program_message)
create_program = HandlerRaw("CreateProgram", [create_code_listing], pogram_sourse_message, pogram_id_message)
update_program_source = HandlerRaw("UpdateProgramSource", [update_code_listing], update_program_sourse_request_message, void_message)
run_program = HandlerRaw("RunProgram", [run_go_code], pogram_id_message, void_message)
lint_program = HandlerRaw("LintProgram", [lint_go_code], pogram_id_message, void_message)
print_all([create_program, get_program_by_id, update_program_source, run_program, lint_program], "Sandbox service")

database_stat = DatabaseRaw([get_program_by_id], [create_program, update_program_source, run_program, lint_program], [create_program])
print_all([database_stat], "Sandbox service databse load")

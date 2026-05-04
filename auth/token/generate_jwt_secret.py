import secrets


def main():
    secret = secrets.token_urlsafe(64)

    print(secret)


if __name__ == "__main__":
    main()

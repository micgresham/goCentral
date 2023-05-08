# goCentral

goCentral is a Go SDK/package that an interface to the HPE/Aruba Networks Central network management system.  This project is at its start and will continue to add functionality over time.  It will strive to offer the same functionality for golang as the pyCentral (https://github.com/aruba/pycentral) SDK offers for Python applications.

If you wish to become a contributor and help to grow this SDK, please contact the author.

goCentral provides a token management system to stores the operational token information in an AES encrypted file as to not require the acquisition of a new token each time an application is executed.  It will handle the refresh of expired tokens when required, but must first be initialized using the Init_DB function.  This function will prompt for the addition of the following information:

  - Central API URL
  - Customer ID
  - Client ID
  - Client secret
  - Access token (this is the current, valid access token)
  - Refresh token

## Functions

------------------------------------------
These functions are NOT public
------------------------------------------


### `func createHash(key string) string`
- Description: Creates a hash value for the given key using the MD5 algorithm.
- Parameters:
  - `key`: The input string for which the hash value needs to be created.
- Returns: The hash value as a hexadecimal string.

### `func Encrypt(data []byte, Passphrase string) []byte`
- Description: Encrypts the provided data using the given passphrase.
- Parameters:
  - `data`: The data to be encrypted.
  - `Passphrase`: The passphrase used for encryption.
- Returns: The encrypted data as a byte slice.

### `func Decrypt(data []byte, Passphrase string) []byte`
- Description: Decrypts the provided data using the given passphrase.
- Parameters:
  - `data`: The data to be decrypted.
  - `Passphrase`: The passphrase used for decryption.
- Returns: The decrypted data as a byte slice.

------------------------------------------
These functions are ARE public
------------------------------------------

### `func RefreshApiToken(central_info Central_struct) (int, string, string, int64)`
- Description: Refreshes the API token for the provided `Central_struct`.
- Parameters:
  - `central_info`: The `Central_struct` containing the necessary information for refreshing the token.
- Returns:
  - The HTTP status code.
  - The new access token.
  - The new refresh token.
  - The expiration time in seconds.

### `func Test_central(central_info Central_struct) (int, string, string)`
- Description: Tests the connectivity to the central API using the provided `Central_struct`.
- Parameters:
  - `central_info`: The `Central_struct` containing the necessary information for testing the API.
- Returns:
  - The HTTP status code.
  - The access token.
  - The refresh token.

### `func Read_DB(filename string) Central_struct`
- Description: Reads the central information from the YAML file specified by `filename` and returns it as a `Central_struct`.
- Parameters:
  - `filename`: The path to the YAML file containing the central information.
- Returns: The `Central_struct` containing the central information.

### `func Write_DB(filename string, central_info_global Central_struct) int`
- Description: Writes the provided `Central_struct` to the YAML file specified by `filename`.
- Parameters:
  - `filename`: The path to the YAML file where the central information will be written.
  - `central_info_global`: The `Central_struct` containing the central information to be written.
- Returns: An integer indicating the success or failure of the write operation (0 for success, 1 for failure).

### `func Init_DB(filename string)`
- Description: Initializes the central database by prompting the user to enter the required information and writing it to the specified YAML file.
- Parameters:
  - `filename`: The path to the YAML file where the central information will be written.



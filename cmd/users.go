package cmd

import (
	"fmt"
	"log"

	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/helpers/prettyprint"
	"github.com/spf13/cobra"
)

var (
	// UserID is the id of a user
	UserID int

	// UserType is the usertype of the user
	UserType string

	// InitialPassword user password
	InitialPassword string

	// AuthenticationMethod user authentication method
	AuthenticationMethod []string

	// UnauthorizedInterfaces  unauthorized interfaces
	UnauthorizedInterfaces []string

	// ExpiryDate when user will expire in EPOCH
	ExpiryDate int

	// VaultAuthorization vault authorization
	VaultAuthorization []string

	// EnableUser if user is enabled
	EnableUser bool

	// ChangePasswordOnLogon if user is prompted to change password on logon
	ChangePasswordOnLogon bool

	// PasswordNeverExpires if user's password will never expire
	PasswordNeverExpires bool

	// DistinguishedName disguished name of user
	DistinguishedName string

	// BusinessAddress of user
	BusinessAddress string

	// Internet info of user
	Internet string

	// Phones of user
	Phones string

	// PersonalDetails of user
	PersonalDetails string
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "User actions for PAS REST API",
	Long: `All users actions that can be taken via PAS REST API.
	
	Example Usage:
	Unsuspend a User: cybr users unsuspend -u userName`,
	Aliases: []string{"user"},
}

var unsuspendUserCmd = &cobra.Command{
	Use:   "unsuspend",
	Short: "Unsuspend a specific user",
	Long: `Activates a suspended user. It does not activate an inactive user.
	
	Example Usage:
	$ cybr users unsuspend --id 9`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.UnsuspendUser(UserID)
		if err != nil {
			log.Fatalf("Failed to unsuspend user with id '%d'. %s", UserID, err)
			return
		}

		fmt.Printf("Successfully unsuspended user with id '%d'\n", UserID)
	},
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List cyberark PAS users",
	Long: `Lists cyberark PAS users.
	
	Example Usage:
	$ cybr users list --search userName --filter userType`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		query := &queries.ListUsers{
			Search: Search,
			Filter: Filter,
		}

		users, err := client.ListUsers(query)
		if err != nil {
			log.Fatalf("Failed to list users. %s", err)
			return
		}

		prettyprint.PrintJSON(users)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cyberark PAS user",
	Long: `Delete a cyberark PAS user.
	
	Example Usage:
	$ cybr users delete --id 9`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		err = client.DeleteUser(UserID)
		if err != nil {
			log.Fatalf("Failed to delete user with id '%d'. %s", UserID, err)
			return
		}

		fmt.Printf("Succesfully deleted user with id '%d'\n", UserID)
	},
}

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a User to PAS",
	Long: `This method adds a new user to the PAS Vault. To run this service one must have 'Add User' and 'Update Users' permissions on the PAS Vault.

	Example Usage:
	$ cybr users add --username userName \
	  --user-type EPVUser \
	  --initial-password initialPassword \
	  --authentication-method AuthTypePass \
	  --location "\\" \
	  --unauthorized-interface PSM,PSMP \
	  --vault-authorization AddSafes,AuditUsers\
	  --enableUser
	  --change-password-on-next-logon
	  --password-never-expires
	  --disguished-name username@Cyberark
	  --description "This is user userName"
	  --internet "homeEmail=userName@Cyberark.com"
	  --personal-details "firtName=me,lastName=lastme"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pasapi.GetConfig()
		if err != nil {
			log.Fatalf("Failed to read configuration file. %s", err)
			return
		}

		businessAddress, err := keyValueStringToMap(BusinessAddress)
		if err != nil {
			log.Fatalf("Failed to parse 'business-address'. %s", err)
		}

		phones, err := keyValueStringToMap(Phones)
		if err != nil {
			log.Fatalf("Failed to parse 'phones'. %s", err)
		}

		personalDetails, err := keyValueStringToMap(PersonalDetails)
		if err != nil {
			log.Fatalf("Failed to parse 'personal-details'. %s", err)
		}

		internet, err := keyValueStringToMap(Internet)
		if err != nil {
			log.Fatalf("Failed to parse 'internet'. %s", err)
		}

		user := requests.AddUser{
			Username:               Username,
			UserType:               UserType,
			InitialPassword:        InitialPassword,
			AuthenticationMethod:   AuthenticationMethod,
			Location:               Location,
			UnAuthorizedInterfaces: UnauthorizedInterfaces,
			ExpiryDate:             ExpiryDate,
			VaultAuthorization:     VaultAuthorization,
			EnableUser:             EnableUser,
			ChangePassOnNextLogon:  ChangePasswordOnLogon,
			PasswordNeverExpires:   PasswordNeverExpires,
			DistinguishedName:      DistinguishedName,
			Description:            Description,
			BusinessAddress:        businessAddress,
			Phones:                 phones,
			PersonalDetails:        personalDetails,
			Internet:               internet,
		}

		response, err := client.AddUser(user)
		if err != nil {
			log.Fatalf("Failed to unsuspend user '%s'. %s", Username, err)
			return
		}

		prettyprint.PrintJSON(response)
	},
}

func init() {
	// unsuspend
	unsuspendUserCmd.Flags().IntVarP(&UserID, "id", "i", 0, "The ID of the user you wish to unsuspend")
	unsuspendUserCmd.MarkFlagRequired("id")

	// list
	listUsersCmd.Flags().StringVarP(&Search, "search", "s", "", "Search for the username, first name or last name of a user")
	listUsersCmd.Flags().StringVarP(&Filter, "filter", "f", "", "Filter on userType or componentUser")

	// delete
	deleteUserCmd.Flags().IntVarP(&UserID, "id", "i", 0, "The ID of the user you wish to delete")
	deleteUserCmd.MarkFlagRequired("id")

	// add
	addUserCmd.Flags().StringVarP(&Username, "username", "u", "", "The username of the user being created")
	addUserCmd.MarkFlagRequired("username")
	addUserCmd.Flags().StringVarP(&Description, "description", "d", "", "The user's notes and comments")
	addUserCmd.Flags().StringVarP(&UserType, "user-type", "t", "EPVUser", "The PAS user type")
	addUserCmd.Flags().StringVarP(&InitialPassword, "initial-password", "p", "", "Initial user password")
	addUserCmd.Flags().StringSliceVarP(&AuthenticationMethod, "authentication-method", "a", []string{"AuthTypePass"}, "User authentication method. Support values: Cyberark, LDAP, Radius")
	addUserCmd.Flags().StringVarP(&Location, "location", "l", "\\", "The location in the Vault where the user will be created")
	addUserCmd.Flags().StringSliceVarP(&UnauthorizedInterfaces, "unauthorized-interfaces", "i", []string{}, "The CyberArk interfaces that this user is not authorized to use")
	addUserCmd.Flags().IntVarP(&ExpiryDate, "expiry-date", "e", 0, "The EPOCH time in which this user expires")
	addUserCmd.Flags().StringSliceVarP(&VaultAuthorization, "vault-authorization", "v", []string{}, "To apply specific authorizations to a user, the user who runs this API must have the same authorizations")
	addUserCmd.Flags().BoolVar(&EnableUser, "enable-user", false, "Whether the user will be enabled upon creation")
	addUserCmd.Flags().BoolVar(&ChangePasswordOnLogon, "change-password-on-logon", false, "Whether or not the user must change their password from the second log on onward")
	addUserCmd.Flags().BoolVar(&PasswordNeverExpires, "password-never-expires", false, "Whether the user’s password will not expire unless they decide to change it")
	addUserCmd.Flags().StringVar(&DistinguishedName, "distinguished-name", "", "The user’s distinguished name. The usage is for PKI authentication, this will match the certificate Subject Name or domain name.")
	addUserCmd.Flags().StringVar(&BusinessAddress, "business-address", "", "The user’s postal address. e.g workCity=Newton,workState=MA")
	addUserCmd.Flags().StringVar(&Internet, "internet", "", "The user's email addresses. e.g homePage=Cyberark.com,homeEmail=user@gmail.com")
	addUserCmd.Flags().StringVar(&Phones, "phones", "", "The user's phone numbers. e.g homeNumber=555123456,businessNumber=555456789")
	addUserCmd.Flags().StringVar(&PersonalDetails, "personal-details", "", "The user's personal details. e.g. street=Dizzengof 56,city=Tel Aviv")

	usersCmd.AddCommand(unsuspendUserCmd)
	usersCmd.AddCommand(listUsersCmd)
	usersCmd.AddCommand(deleteUserCmd)
	usersCmd.AddCommand(addUserCmd)
	rootCmd.AddCommand(usersCmd)
}

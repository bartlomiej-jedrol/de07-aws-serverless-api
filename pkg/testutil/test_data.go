package testutil

import (
	"fmt"

	"github.com/bartlomiej-jedrol/de07-aws-serverless-api/pkg/models"
)

var (
	ValidUser1 = models.User{
		Email:     "bartlomiej.jedrol@gmail.com",
		FirstName: "Bartlomiej",
		LastName:  "Jedrol",
		Age:       37,
	}

	ValidUser2 = models.User{
		Email:     "jedrol.natalia@gmail.com",
		FirstName: "Natalia",
		LastName:  "Jedrol",
		Age:       33,
	}

	InvalidUser1 = models.User{
		Email:     "test.test@gmail.com",
		FirstName: "test",
		LastName:  "test",
		Age:       1,
	}

	ValidUser string = fmt.Sprintf(`{"email":"%v","firstName":"%v","lastName":"%v","age":%v}`,
		ValidUser1.Email, ValidUser1.FirstName, ValidUser1.LastName, ValidUser1.Age)

	InvalidUser string = fmt.Sprintf(`{"email":"%v","firstName":"%v","lastName":"%v","age":%v}`,
		InvalidUser1.Email, InvalidUser1.FirstName, InvalidUser1.LastName, InvalidUser1.Age)

	UserEmptyEmail string = fmt.Sprintf(`{"email":"","firstName":"%v","lastName":"%v","age":%v}`,
		ValidUser1.FirstName, ValidUser1.LastName, ValidUser1.Age)

	EmptyUser string = `{}`

	InvalidJSON string = fmt.Sprintf(`{"email":""%v","firstName":"%v","lastName":"%v","age":%v}`,
		ValidUser1.Email, ValidUser1.FirstName, ValidUser1.LastName, ValidUser1.Age)

	ValidUsers string = fmt.Sprintf(`[{"email":"%v","firstName":"%v","lastName":"%v","age":%v},{"email":"%v","firstName":"%v","lastName":"%v","age":%v}]`,
		ValidUser2.Email, ValidUser2.FirstName, ValidUser2.LastName, ValidUser2.Age, ValidUser1.Email, ValidUser1.FirstName, ValidUser1.LastName, ValidUser1.Age)

	ValidQueQueryStringParameters   = map[string]string{"email": ValidUser1.Email}
	InvalidQueQueryStringParameters = map[string]string{"email": InvalidUser1.Email}
)

package handlers

import (
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"gorm.io/gorm"
)

func VerifyNewUser(db *gorm.DB, verification models.Verification) bool {

	ok1 := UpdateField(db, models.User{}, "id", verification.EntityID, "verified", true)
	ok2 := UpdateField(db, models.Passkey{}, "user_id", verification.EntityID, "verified", true)
	ok3 := UpdateField(db, models.AccessToken{}, "user_id", verification.EntityID, "disabled", false)

	verification.Status = models.StatusSuccess

	res := db.Save(&verification)

	return ok1 && ok2 && ok3 && res.RowsAffected != 0 && res.Error == nil
}

func VerifyNewPasskey(db *gorm.DB, verification models.Verification) bool {

	ok := UpdateField(db, models.Passkey{}, "id", verification.EntityID, "verified", true)

	verification.Status = models.StatusSuccess

	res := db.Save(&verification)

	return ok && res.RowsAffected != 0 && res.Error == nil
}

func DeleteUser(db *gorm.DB, verification models.Verification) bool {

	ok := DeleteInDatabaseById(db, "id", verification.EntityID, &[]models.User{})
	DeleteInDatabaseById(db, "user_id", verification.EntityID, &[]models.Passkey{})
	DeleteInDatabaseById(db, "user_id", verification.EntityID, &[]models.Challenge{})
	DeleteInDatabaseById(db, "user_id", verification.EntityID, &[]models.AccessToken{})
	DeleteInDatabaseById(db, "user_id", verification.EntityID, &[]models.AccessLog{})
	DeleteInDatabaseById(db, "user_id", verification.EntityID, &[]models.Verification{})

	// INFO: PRIVATE KEY: Uncomment if we need to Store Private Keys
	// DeleteInDatabaseById(c, db, "user_id", verification.EntityID, &[]models.PasskeyPrivateKey{})

	// INFO Dont Update Verificaion status as row will be already deleted
	return ok
}

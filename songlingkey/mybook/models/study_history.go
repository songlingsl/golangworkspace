package models

type StudyHistory struct {
	StudyHistorId int `gorm:"PRIMARY_KEY:book_Id"`
	UserId        int
	BookId        string
	NickName      string
	BookName      string
}

func SaveHistory(history *StudyHistory) error {

	if err := db.Create(history).Error; err != nil {
		return err
	}
	return nil

}

func UpdateHistory(history *StudyHistory) int64 {
	return db.Exec("UPDATE study_history SET counts=counts+1 WHERE user_Id = ? and book_id=?", history.UserId, history.BookId).RowsAffected
	//UPDATE study_history SET counts=counts+1 WHERE study_history_id = 13
	//return db.Model(&StudyHistory{}).Where("user_Id = ? and book_id=?", history.UserId, history.BookId).Updates(StudyHistory{UserId: history.UserId}).RowsAffected

}

package service

import (
	"github.com/s.usynin/testing/go-server/internal/models"
	"github.com/s.usynin/testing/go-server/internal/repository"
)

type PostService struct {
	postRepo     *repository.PostRepository
	commentRepo  *repository.CommentRepository
	categoryRepo *repository.CategoryRepository
	likeRepo     *repository.LikeRepository
}

func NewPostService(
	postRepo *repository.PostRepository,
	commentRepo *repository.CommentRepository,
	categoryRepo *repository.CategoryRepository,
	likeRepo *repository.LikeRepository,
) *PostService {
	return &PostService{
		postRepo:     postRepo,
		commentRepo:  commentRepo,
		categoryRepo: categoryRepo,
		likeRepo:     likeRepo,
	}
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	posts, err := s.postRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Загружаем категории для постов
	for i := range posts {
		category, err := s.categoryRepo.GetByID(posts[i].CategoryID)
		if err == nil {
			posts[i].Category = category
		}
	}

	return posts, nil
}

func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Загружаем категорию
	category, err := s.categoryRepo.GetByID(post.CategoryID)
	if err == nil {
		post.Category = category
	}

	// Загружаем комментарии
	comments, err := s.commentRepo.GetByPostID(id)
	if err == nil {
		post.Comments = comments
	}

	return post, nil
}

func (s *PostService) CreatePost(title, content string, categoryID int) (*models.Post, error) {
	id, err := s.postRepo.Create(title, content, categoryID)
	if err != nil {
		return nil, err
	}

	return s.GetPostByID(int(id))
}

func (s *PostService) DeletePost(id int) error {
	return s.postRepo.Delete(id)
}

func (s *PostService) GetPostsByCategory(categoryID int) ([]models.Post, error) {
	return s.postRepo.GetByCategoryID(categoryID)
}

func (s *PostService) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *PostService) AddComment(postID int, author, content string) (*models.Comment, error) {
	id, err := s.commentRepo.Create(postID, author, content)
	if err != nil {
		return nil, err
	}

	return s.commentRepo.GetByID(int(id))
}

func (s *PostService) AddLike(postID int) (int64, error) {
	return s.likeRepo.Add(postID)
}

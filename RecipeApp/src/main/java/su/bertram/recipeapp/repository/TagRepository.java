package su.bertram.recipeapp.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import su.bertram.recipeapp.model.Tag;

public interface TagRepository extends JpaRepository<Tag, Long> {

}

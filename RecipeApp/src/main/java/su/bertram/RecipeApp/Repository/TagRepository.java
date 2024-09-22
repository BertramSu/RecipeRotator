package su.bertram.RecipeApp.Repository;

import org.springframework.data.jpa.repository.JpaRepository;
import su.bertram.RecipeApp.Model.Tag;

public interface TagRepository extends JpaRepository<Tag, Long> {

}
